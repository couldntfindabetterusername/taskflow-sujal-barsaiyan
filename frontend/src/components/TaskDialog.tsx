import { useState, useEffect, FormEvent } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  Alert,
  CircularProgress,
  Box,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Grid,
} from '@mui/material';
import type { Task, TaskCreateRequest, TaskUpdateRequest, TaskStatus, TaskPriority } from '../types';

interface TaskDialogProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: TaskCreateRequest | TaskUpdateRequest) => Promise<void>;
  task?: Task | null; // If provided, we're editing; otherwise creating
  isLoading?: boolean;
  error?: string | null;
}

const statusOptions: { value: TaskStatus; label: string }[] = [
  { value: 'todo', label: 'To Do' },
  { value: 'in_progress', label: 'In Progress' },
  { value: 'done', label: 'Done' },
];

const priorityOptions: { value: TaskPriority; label: string }[] = [
  { value: 'low', label: 'Low' },
  { value: 'medium', label: 'Medium' },
  { value: 'high', label: 'High' },
];

export default function TaskDialog({
  open,
  onClose,
  onSubmit,
  task,
  isLoading = false,
  error,
}: TaskDialogProps) {
  const isEditing = !!task;

  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [status, setStatus] = useState<TaskStatus>('todo');
  const [priority, setPriority] = useState<TaskPriority>('medium');
  const [dueDate, setDueDate] = useState('');
  const [validationErrors, setValidationErrors] = useState<{
    title?: string;
    dueDate?: string;
  }>({});

  // Reset form when dialog opens or task changes
  useEffect(() => {
    if (open) {
      if (task) {
        setTitle(task.title);
        setDescription(task.description || '');
        setStatus(task.status);
        setPriority(task.priority);
        setDueDate(task.due_date ? task.due_date.split('T')[0] : '');
      } else {
        setTitle('');
        setDescription('');
        setStatus('todo');
        setPriority('medium');
        setDueDate('');
      }
      setValidationErrors({});
    }
  }, [open, task]);

  const validateForm = (): boolean => {
    const errors: { title?: string; dueDate?: string } = {};

    if (!title.trim()) {
      errors.title = 'Task title is required';
    } else if (title.trim().length < 2) {
      errors.title = 'Task title must be at least 2 characters';
    } else if (title.trim().length > 200) {
      errors.title = 'Task title must be less than 200 characters';
    }

    if (dueDate) {
      const selectedDate = new Date(dueDate);
      const today = new Date();
      today.setHours(0, 0, 0, 0);

      // Only validate future date for new tasks, not edits
      if (!isEditing && selectedDate < today) {
        errors.dueDate = 'Due date must be in the future';
      }
    }

    setValidationErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    const data: TaskCreateRequest | TaskUpdateRequest = {
      title: title.trim(),
      description: description.trim() || undefined,
      status,
      priority,
      due_date: dueDate || undefined,
    };

    await onSubmit(data);
  };

  const handleClose = () => {
    if (!isLoading) {
      onClose();
    }
  };

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      maxWidth="sm"
      fullWidth
      aria-labelledby="task-dialog-title"
    >
      <DialogTitle id="task-dialog-title">
        {isEditing ? 'Edit Task' : 'Create New Task'}
      </DialogTitle>
      <Box component="form" onSubmit={handleSubmit}>
        <DialogContent>
          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}

          {/* Title */}
          <TextField
            autoFocus
            fullWidth
            label="Task Title"
            value={title}
            onChange={(e) => {
              setTitle(e.target.value);
              if (validationErrors.title) {
                setValidationErrors((prev) => ({ ...prev, title: undefined }));
              }
            }}
            error={!!validationErrors.title}
            helperText={validationErrors.title}
            margin="normal"
            disabled={isLoading}
            required
          />

          {/* Description */}
          <TextField
            fullWidth
            label="Description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            margin="normal"
            multiline
            rows={3}
            disabled={isLoading}
            placeholder="Optional: Add a description for your task"
          />

          {/* Status and Priority Row */}
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth size="small">
                <InputLabel id="status-label">Status</InputLabel>
                <Select
                  labelId="status-label"
                  value={status}
                  label="Status"
                  onChange={(e) => setStatus(e.target.value as TaskStatus)}
                  disabled={isLoading}
                >
                  {statusOptions.map((option) => (
                    <MenuItem key={option.value} value={option.value}>
                      {option.label}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth size="small">
                <InputLabel id="priority-label">Priority</InputLabel>
                <Select
                  labelId="priority-label"
                  value={priority}
                  label="Priority"
                  onChange={(e) => setPriority(e.target.value as TaskPriority)}
                  disabled={isLoading}
                >
                  {priorityOptions.map((option) => (
                    <MenuItem key={option.value} value={option.value}>
                      {option.label}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
          </Grid>

          {/* Due Date */}
          <TextField
            fullWidth
            label="Due Date"
            type="date"
            value={dueDate}
            onChange={(e) => {
              setDueDate(e.target.value);
              if (validationErrors.dueDate) {
                setValidationErrors((prev) => ({ ...prev, dueDate: undefined }));
              }
            }}
            error={!!validationErrors.dueDate}
            helperText={validationErrors.dueDate || 'Optional'}
            margin="normal"
            disabled={isLoading}
            InputLabelProps={{
              shrink: true,
            }}
          />
        </DialogContent>

        <DialogActions sx={{ px: 3, pb: 2 }}>
          <Button onClick={handleClose} disabled={isLoading}>
            Cancel
          </Button>
          <Button
            type="submit"
            variant="contained"
            disabled={isLoading}
            startIcon={isLoading ? <CircularProgress size={16} /> : null}
          >
            {isLoading ? (isEditing ? 'Saving...' : 'Creating...') : (isEditing ? 'Save Changes' : 'Create Task')}
          </Button>
        </DialogActions>
      </Box>
    </Dialog>
  );
}
