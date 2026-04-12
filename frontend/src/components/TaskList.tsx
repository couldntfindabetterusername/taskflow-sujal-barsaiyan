import { Box, Typography, Paper } from '@mui/material';
import AssignmentIcon from '@mui/icons-material/Assignment';
import TaskItem from './TaskItem';
import type { Task, TaskStatus } from '../types';

interface TaskListProps {
  tasks: Task[];
  onStatusChange: (taskId: string, newStatus: TaskStatus) => void;
  onEdit: (task: Task) => void;
  onDelete: (task: Task) => void;
  updatingTaskId?: string | null;
  emptyMessage?: string;
}

export default function TaskList({
  tasks,
  onStatusChange,
  onEdit,
  onDelete,
  updatingTaskId,
  emptyMessage = 'No tasks found',
}: TaskListProps) {
  if (tasks.length === 0) {
    return (
      <Paper
        sx={{
          p: 4,
          textAlign: 'center',
          backgroundColor: 'background.paper',
        }}
      >
        <AssignmentIcon
          sx={{
            fontSize: 60,
            color: 'text.disabled',
            mb: 2,
          }}
        />
        <Typography variant="h6" color="text.secondary" gutterBottom>
          {emptyMessage}
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Create a new task to get started or adjust your filters.
        </Typography>
      </Paper>
    );
  }

  return (
    <Box>
      {tasks.map((task) => (
        <TaskItem
          key={task.id}
          task={task}
          onStatusChange={onStatusChange}
          onEdit={onEdit}
          onDelete={onDelete}
          isUpdating={updatingTaskId === task.id}
        />
      ))}
    </Box>
  );
}
