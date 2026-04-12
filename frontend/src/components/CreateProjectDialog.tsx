import { useState, FormEvent } from 'react';
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
} from '@mui/material';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { projectsApi, getErrorMessage } from '../services/api';
import type { ProjectCreateRequest } from '../types';

interface CreateProjectDialogProps {
  open: boolean;
  onClose: () => void;
}

export default function CreateProjectDialog({ open, onClose }: CreateProjectDialogProps) {
  const queryClient = useQueryClient();

  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [validationErrors, setValidationErrors] = useState<{
    name?: string;
  }>({});
  const [apiError, setApiError] = useState<string | null>(null);

  const createMutation = useMutation({
    mutationFn: (data: ProjectCreateRequest) => projectsApi.create(data),
    onSuccess: () => {
      // Invalidate projects query to refetch the list
      queryClient.invalidateQueries({ queryKey: ['projects'] });
      handleClose();
    },
    onError: (error) => {
      setApiError(getErrorMessage(error));
    },
  });

  const validateForm = (): boolean => {
    const errors: { name?: string } = {};

    if (!name.trim()) {
      errors.name = 'Project name is required';
    } else if (name.trim().length < 2) {
      errors.name = 'Project name must be at least 2 characters';
    } else if (name.trim().length > 100) {
      errors.name = 'Project name must be less than 100 characters';
    }

    setValidationErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setApiError(null);

    if (!validateForm()) {
      return;
    }

    createMutation.mutate({
      name: name.trim(),
      description: description.trim() || undefined,
    });
  };

  const handleClose = () => {
    // Reset form state
    setName('');
    setDescription('');
    setValidationErrors({});
    setApiError(null);
    createMutation.reset();
    onClose();
  };

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      maxWidth="sm"
      fullWidth
      aria-labelledby="create-project-dialog-title"
    >
      <DialogTitle id="create-project-dialog-title">Create New Project</DialogTitle>
      <Box component="form" onSubmit={handleSubmit}>
        <DialogContent>
          {apiError && (
            <Alert severity="error" sx={{ mb: 2 }} onClose={() => setApiError(null)}>
              {apiError}
            </Alert>
          )}

          <TextField
            autoFocus
            fullWidth
            label="Project Name"
            value={name}
            onChange={(e) => {
              setName(e.target.value);
              if (validationErrors.name) {
                setValidationErrors((prev) => ({ ...prev, name: undefined }));
              }
            }}
            error={!!validationErrors.name}
            helperText={validationErrors.name}
            margin="normal"
            disabled={createMutation.isPending}
            required
          />

          <TextField
            fullWidth
            label="Description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            margin="normal"
            multiline
            rows={3}
            disabled={createMutation.isPending}
            placeholder="Optional: Add a description for your project"
          />
        </DialogContent>

        <DialogActions sx={{ px: 3, pb: 2 }}>
          <Button onClick={handleClose} disabled={createMutation.isPending}>
            Cancel
          </Button>
          <Button
            type="submit"
            variant="contained"
            disabled={createMutation.isPending}
            startIcon={createMutation.isPending ? <CircularProgress size={16} /> : null}
          >
            {createMutation.isPending ? 'Creating...' : 'Create Project'}
          </Button>
        </DialogActions>
      </Box>
    </Dialog>
  );
}
