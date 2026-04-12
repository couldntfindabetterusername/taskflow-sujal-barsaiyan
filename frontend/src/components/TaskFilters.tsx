import {
  Box,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  SelectChangeEvent,
} from '@mui/material';
import type { TaskStatus } from '../types';

interface TaskFiltersProps {
  statusFilter: TaskStatus | 'all';
  assigneeFilter: string;
  onStatusChange: (status: TaskStatus | 'all') => void;
  onAssigneeChange: (assignee: string) => void;
  currentUserId: string;
}

export default function TaskFilters({
  statusFilter,
  assigneeFilter,
  onStatusChange,
  onAssigneeChange,
  currentUserId,
}: TaskFiltersProps) {
  const handleStatusChange = (event: SelectChangeEvent<string>) => {
    onStatusChange(event.target.value as TaskStatus | 'all');
  };

  const handleAssigneeChange = (event: SelectChangeEvent<string>) => {
    onAssigneeChange(event.target.value);
  };

  return (
    <Box
      sx={{
        display: 'flex',
        gap: 2,
        flexWrap: 'wrap',
        mb: 3,
      }}
    >
      {/* Status Filter */}
      <FormControl size="small" sx={{ minWidth: 150 }}>
        <InputLabel id="status-filter-label">Status</InputLabel>
        <Select
          labelId="status-filter-label"
          id="status-filter"
          value={statusFilter}
          label="Status"
          onChange={handleStatusChange}
        >
          <MenuItem value="all">All Statuses</MenuItem>
          <MenuItem value="todo">To Do</MenuItem>
          <MenuItem value="in_progress">In Progress</MenuItem>
          <MenuItem value="done">Done</MenuItem>
        </Select>
      </FormControl>

      {/* Assignee Filter */}
      <FormControl size="small" sx={{ minWidth: 150 }}>
        <InputLabel id="assignee-filter-label">Assignee</InputLabel>
        <Select
          labelId="assignee-filter-label"
          id="assignee-filter"
          value={assigneeFilter}
          label="Assignee"
          onChange={handleAssigneeChange}
        >
          <MenuItem value="all">All Assignees</MenuItem>
          <MenuItem value={currentUserId}>Assigned to Me</MenuItem>
          <MenuItem value="unassigned">Unassigned</MenuItem>
        </Select>
      </FormControl>
    </Box>
  );
}
