import {
  Card,
  CardContent,
  Box,
  Typography,
  Chip,
  IconButton,
  Select,
  MenuItem,
  FormControl,
  SelectChangeEvent,
  Tooltip,
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import PersonIcon from '@mui/icons-material/Person';
import CalendarTodayIcon from '@mui/icons-material/CalendarToday';
import type { Task, TaskStatus, TaskPriority } from '../types';

interface TaskItemProps {
  task: Task;
  onStatusChange: (taskId: string, newStatus: TaskStatus) => void;
  onEdit: (task: Task) => void;
  onDelete: (task: Task) => void;
  isUpdating?: boolean;
}

const statusColors: Record<TaskStatus, 'default' | 'primary' | 'success'> = {
  todo: 'default',
  in_progress: 'primary',
  done: 'success',
};

const statusLabels: Record<TaskStatus, string> = {
  todo: 'To Do',
  in_progress: 'In Progress',
  done: 'Done',
};

const priorityColors: Record<TaskPriority, 'default' | 'warning' | 'error'> = {
  low: 'default',
  medium: 'warning',
  high: 'error',
};

const priorityLabels: Record<TaskPriority, string> = {
  low: 'Low',
  medium: 'Medium',
  high: 'High',
};

export default function TaskItem({
  task,
  onStatusChange,
  onEdit,
  onDelete,
  isUpdating = false,
}: TaskItemProps) {
  const handleStatusChange = (event: SelectChangeEvent<TaskStatus>) => {
    event.stopPropagation();
    onStatusChange(task.id, event.target.value as TaskStatus);
  };

  const formatDate = (dateString: string | null): string => {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  };

  const isOverdue = (): boolean => {
    if (!task.due_date || task.status === 'done') return false;
    return new Date(task.due_date) < new Date();
  };

  return (
    <Card
      sx={{
        mb: 2,
        opacity: isUpdating ? 0.7 : 1,
        transition: 'opacity 0.2s, box-shadow 0.2s',
        '&:hover': {
          boxShadow: 3,
        },
        borderLeft: 4,
        borderLeftColor:
          task.status === 'done'
            ? 'success.main'
            : task.status === 'in_progress'
            ? 'primary.main'
            : 'grey.300',
      }}
    >
      <CardContent sx={{ py: 2, '&:last-child': { pb: 2 } }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
          {/* Main Content */}
          <Box sx={{ flexGrow: 1, mr: 2 }}>
            {/* Title */}
            <Typography
              variant="subtitle1"
              component="h3"
              sx={{
                fontWeight: 600,
                textDecoration: task.status === 'done' ? 'line-through' : 'none',
                color: task.status === 'done' ? 'text.secondary' : 'text.primary',
                mb: 0.5,
              }}
            >
              {task.title}
            </Typography>

            {/* Description */}
            {task.description && (
              <Typography
                variant="body2"
                color="text.secondary"
                sx={{
                  mb: 1.5,
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                  display: '-webkit-box',
                  WebkitLineClamp: 2,
                  WebkitBoxOrient: 'vertical',
                }}
              >
                {task.description}
              </Typography>
            )}

            {/* Metadata Row */}
            <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, alignItems: 'center' }}>
              {/* Priority Chip */}
              <Chip
                label={priorityLabels[task.priority]}
                size="small"
                color={priorityColors[task.priority]}
                variant="outlined"
                sx={{ fontSize: '0.75rem' }}
              />

              {/* Assignee */}
              {task.assignee && (
                <Chip
                  icon={<PersonIcon sx={{ fontSize: '0.875rem' }} />}
                  label={task.assignee.name}
                  size="small"
                  variant="outlined"
                  sx={{ fontSize: '0.75rem' }}
                />
              )}

              {/* Due Date */}
              {task.due_date && (
                <Chip
                  icon={<CalendarTodayIcon sx={{ fontSize: '0.875rem' }} />}
                  label={formatDate(task.due_date)}
                  size="small"
                  variant="outlined"
                  color={isOverdue() ? 'error' : 'default'}
                  sx={{ fontSize: '0.75rem' }}
                />
              )}
            </Box>
          </Box>

          {/* Actions */}
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            {/* Status Dropdown */}
            <FormControl size="small" sx={{ minWidth: 120 }}>
              <Select
                value={task.status}
                onChange={handleStatusChange}
                disabled={isUpdating}
                sx={{
                  '& .MuiSelect-select': {
                    py: 0.75,
                    fontSize: '0.875rem',
                  },
                }}
              >
                <MenuItem value="todo">
                  <Chip
                    label={statusLabels.todo}
                    size="small"
                    color={statusColors.todo}
                    sx={{ fontSize: '0.75rem' }}
                  />
                </MenuItem>
                <MenuItem value="in_progress">
                  <Chip
                    label={statusLabels.in_progress}
                    size="small"
                    color={statusColors.in_progress}
                    sx={{ fontSize: '0.75rem' }}
                  />
                </MenuItem>
                <MenuItem value="done">
                  <Chip
                    label={statusLabels.done}
                    size="small"
                    color={statusColors.done}
                    sx={{ fontSize: '0.75rem' }}
                  />
                </MenuItem>
              </Select>
            </FormControl>

            {/* Edit Button */}
            <Tooltip title="Edit task">
              <IconButton
                size="small"
                onClick={() => onEdit(task)}
                disabled={isUpdating}
                aria-label="Edit task"
              >
                <EditIcon fontSize="small" />
              </IconButton>
            </Tooltip>

            {/* Delete Button */}
            <Tooltip title="Delete task">
              <IconButton
                size="small"
                onClick={() => onDelete(task)}
                disabled={isUpdating}
                color="error"
                aria-label="Delete task"
              >
                <DeleteIcon fontSize="small" />
              </IconButton>
            </Tooltip>
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
}
