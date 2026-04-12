import { useState, useMemo } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box,
  Typography,
  Button,
  Alert,
  CircularProgress,
  IconButton,
  Breadcrumbs,
  Link,
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { projectsApi, tasksApi, getErrorMessage } from '../services/api';
import { useAuth } from '../context/AuthContext';
import TaskFilters from '../components/TaskFilters';
import TaskList from '../components/TaskList';
import TaskDialog from '../components/TaskDialog';
import DeleteConfirmDialog from '../components/DeleteConfirmDialog';
import type { Task, TaskStatus, TaskCreateRequest, TaskUpdateRequest } from '../types';

export default function ProjectDetailPage() {
  const { id: projectId } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { user } = useAuth();

  // Filter state
  const [statusFilter, setStatusFilter] = useState<TaskStatus | 'all'>('all');
  const [assigneeFilter, setAssigneeFilter] = useState<string>('all');

  // Dialog state
  const [taskDialogOpen, setTaskDialogOpen] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | null>(null);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [taskToDelete, setTaskToDelete] = useState<Task | null>(null);

  // Track which task is being updated (for optimistic UI)
  const [updatingTaskId, setUpdatingTaskId] = useState<string | null>(null);

  // Fetch project with tasks
  const {
    data: projectData,
    isLoading,
    isError,
    error,
    refetch,
  } = useQuery({
    queryKey: ['project', projectId],
    queryFn: () => projectsApi.get(projectId!),
    enabled: !!projectId,
  });

  // Create task mutation
  const createTaskMutation = useMutation({
    mutationFn: (data: TaskCreateRequest) => tasksApi.create(projectId!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['project', projectId] });
      setTaskDialogOpen(false);
    },
  });

  // Update task mutation with optimistic updates
  const updateTaskMutation = useMutation({
    mutationFn: ({ taskId, data }: { taskId: string; data: TaskUpdateRequest }) =>
      tasksApi.update(taskId, data),
    onMutate: async ({ taskId, data }) => {
      // Cancel outgoing refetches
      await queryClient.cancelQueries({ queryKey: ['project', projectId] });

      // Snapshot previous value
      const previousData = queryClient.getQueryData(['project', projectId]);

      // Optimistically update
      queryClient.setQueryData(['project', projectId], (old: typeof projectData) => {
        if (!old) return old;
        return {
          ...old,
          tasks: old.tasks.map((task: Task) =>
            task.id === taskId ? { ...task, ...data } : task
          ),
        };
      });

      setUpdatingTaskId(taskId);

      return { previousData };
    },
    onError: (_err, _variables, context) => {
      // Rollback on error
      if (context?.previousData) {
        queryClient.setQueryData(['project', projectId], context.previousData);
      }
    },
    onSettled: () => {
      setUpdatingTaskId(null);
      queryClient.invalidateQueries({ queryKey: ['project', projectId] });
    },
  });

  // Delete task mutation
  const deleteTaskMutation = useMutation({
    mutationFn: (taskId: string) => tasksApi.delete(taskId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['project', projectId] });
      setDeleteDialogOpen(false);
      setTaskToDelete(null);
    },
  });

  // Filter tasks
  const filteredTasks = useMemo(() => {
    if (!projectData?.tasks) return [];

    return projectData.tasks.filter((task: Task) => {
      // Status filter
      if (statusFilter !== 'all' && task.status !== statusFilter) {
        return false;
      }

      // Assignee filter
      if (assigneeFilter === 'unassigned' && task.assignee_id !== null) {
        return false;
      }
      if (assigneeFilter !== 'all' && assigneeFilter !== 'unassigned' && task.assignee_id !== assigneeFilter) {
        return false;
      }

      return true;
    });
  }, [projectData?.tasks, statusFilter, assigneeFilter]);

  // Handlers
  const handleStatusChange = (taskId: string, newStatus: TaskStatus) => {
    updateTaskMutation.mutate({ taskId, data: { status: newStatus } });
  };

  const handleOpenCreateDialog = () => {
    setEditingTask(null);
    setTaskDialogOpen(true);
  };

  const handleOpenEditDialog = (task: Task) => {
    setEditingTask(task);
    setTaskDialogOpen(true);
  };

  const handleCloseTaskDialog = () => {
    setTaskDialogOpen(false);
    setEditingTask(null);
    createTaskMutation.reset();
    updateTaskMutation.reset();
  };

  const handleTaskSubmit = async (data: TaskCreateRequest | TaskUpdateRequest) => {
    if (editingTask) {
      await updateTaskMutation.mutateAsync({ taskId: editingTask.id, data });
      setTaskDialogOpen(false);
      setEditingTask(null);
    } else {
      await createTaskMutation.mutateAsync(data as TaskCreateRequest);
    }
  };

  const handleOpenDeleteDialog = (task: Task) => {
    setTaskToDelete(task);
    setDeleteDialogOpen(true);
  };

  const handleCloseDeleteDialog = () => {
    setDeleteDialogOpen(false);
    setTaskToDelete(null);
  };

  const handleConfirmDelete = () => {
    if (taskToDelete) {
      deleteTaskMutation.mutate(taskToDelete.id);
    }
  };

  // Loading state
  if (isLoading) {
    return (
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          minHeight: '50vh',
        }}
      >
        <CircularProgress />
      </Box>
    );
  }

  // Error state
  if (isError) {
    return (
      <Box sx={{ py: 4 }}>
        <Alert
          severity="error"
          action={
            <Button color="inherit" size="small" onClick={() => refetch()}>
              Retry
            </Button>
          }
        >
          {getErrorMessage(error)}
        </Alert>
      </Box>
    );
  }

  if (!projectData) {
    return (
      <Box sx={{ py: 4 }}>
        <Alert severity="error">Project not found</Alert>
      </Box>
    );
  }

  const { project, tasks } = projectData;

  // Get empty message based on filters
  const getEmptyMessage = () => {
    if (tasks.length === 0) {
      return 'No tasks in this project yet';
    }
    if (statusFilter !== 'all' || assigneeFilter !== 'all') {
      return 'No tasks match your filters';
    }
    return 'No tasks found';
  };

  return (
    <Box sx={{ py: 2 }}>
      {/* Breadcrumbs */}
      <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
        <IconButton
          onClick={() => navigate('/projects')}
          sx={{ mr: 1 }}
          aria-label="Back to projects"
        >
          <ArrowBackIcon />
        </IconButton>
        <Breadcrumbs aria-label="breadcrumb">
          <Link
            component="button"
            variant="body1"
            onClick={() => navigate('/projects')}
            underline="hover"
            color="inherit"
            sx={{ cursor: 'pointer' }}
          >
            Projects
          </Link>
          <Typography color="text.primary">{project.name}</Typography>
        </Breadcrumbs>
      </Box>

      {/* Project Header */}
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          {project.name}
        </Typography>
        {project.description && (
          <Typography variant="body1" color="text.secondary">
            {project.description}
          </Typography>
        )}
      </Box>

      {/* Tasks Section */}
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          mb: 2,
        }}
      >
        <Typography variant="h5" component="h2">
          Tasks ({filteredTasks.length})
        </Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={handleOpenCreateDialog}
        >
          New Task
        </Button>
      </Box>

      {/* Filters */}
      <TaskFilters
        statusFilter={statusFilter}
        assigneeFilter={assigneeFilter}
        onStatusChange={setStatusFilter}
        onAssigneeChange={setAssigneeFilter}
        currentUserId={user?.id || ''}
      />

      {/* Task List */}
      <TaskList
        tasks={filteredTasks}
        onStatusChange={handleStatusChange}
        onEdit={handleOpenEditDialog}
        onDelete={handleOpenDeleteDialog}
        updatingTaskId={updatingTaskId}
        emptyMessage={getEmptyMessage()}
      />

      {/* Create/Edit Task Dialog */}
      <TaskDialog
        open={taskDialogOpen}
        onClose={handleCloseTaskDialog}
        onSubmit={handleTaskSubmit}
        task={editingTask}
        isLoading={createTaskMutation.isPending || updateTaskMutation.isPending}
        error={
          createTaskMutation.isError
            ? getErrorMessage(createTaskMutation.error)
            : updateTaskMutation.isError
            ? getErrorMessage(updateTaskMutation.error)
            : null
        }
      />

      {/* Delete Confirmation Dialog */}
      <DeleteConfirmDialog
        open={deleteDialogOpen}
        onClose={handleCloseDeleteDialog}
        onConfirm={handleConfirmDelete}
        title="Delete Task"
        message={`Are you sure you want to delete "${taskToDelete?.title}"? This action cannot be undone.`}
        isLoading={deleteTaskMutation.isPending}
      />
    </Box>
  );
}
