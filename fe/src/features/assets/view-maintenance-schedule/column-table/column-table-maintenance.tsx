import type { ColumnDef } from '@tanstack/react-table'
import { format } from 'date-fns'
import { Badge, Button } from '@/components/ui'
import { Link } from 'react-router-dom'
import type { MaintenanceSchedule } from '../model'
import { Calendar } from 'lucide-react'
import { UpdateMaintenanceSchedule } from '../../update-maintenance-schedule'
import { useState } from 'react'
import { useAppSelector } from '@/hooks'
export const columnTableMaintenance = ({
  onSuccessUpdate,
}: {
  onSuccessUpdate: () => void
}): ColumnDef<MaintenanceSchedule>[] => [
  {
    accessorKey: 'asset.assetName',
    header: 'Asset Name',
    cell: ({ row }) => (
      <div>
        <div className='font-medium'>{row.original.asset.assetName}</div>
        <div className='text-muted-foreground text-xs'>ID: {row.original.asset.id}</div>
      </div>
    ),
  },
  {
    accessorKey: 'startDate',
    header: 'Start Date',
    cell: ({ row }) => format(new Date(row.original.startDate), 'MM/dd/yyyy'),
  },
  {
    accessorKey: 'endDate',
    header: 'End Date',
    cell: ({ row }) => format(new Date(row.original.endDate), 'MM/dd/yyyy'),
  },
  {
    accessorKey: 'asset.status',
    header: 'Status',
    cell: ({ row }) => (
      <Badge
        variant='outline'
        className={
          row.original.asset.status === 'New'
            ? 'border-green-200 bg-green-100 text-green-800'
            : row.original.asset.status === 'In Use'
              ? 'border-blue-200 bg-blue-100 text-blue-800'
              : row.original.asset.status === 'Under Maintenance'
                ? 'border-yellow-200 bg-yellow-100 text-yellow-800'
                : 'border-red-200 bg-red-100 text-red-800'
        }
      >
        {row.original.asset.status}
      </Badge>
    ),
  },
  {
    id: 'actions',
    header: '',
    cell: ({ row }) => {
      const id = row.original.id
      const status = row.original.asset.status
      const [isDialogOpen, setIsDialogOpen] = useState(false)
      const role = useAppSelector((state) => state.auth.user?.role.slug)
      const startDate = row.original.startDate
      const endDate = row.original.endDate
      return (
        <div className='flex items-center gap-2'>
          <UpdateMaintenanceSchedule
            startDate={startDate}
            endDate={endDate}
            id={id.toString()}
            isDialogOpen={isDialogOpen}
            setIsDialogOpen={setIsDialogOpen}
            onSuccessUpdate={onSuccessUpdate}
          />
          <Button
            variant='outline'
            size='sm'
            asChild
          >
            <Link to={`/assets/${row.original.asset.id}`}>View Asset</Link>
          </Button>
          {role !== 'viewer' && status !== 'Under Maintenance' && (
            <Button
              variant='outline'
              size='sm'
              onClick={() => {
                setIsDialogOpen(true)
              }}
            >
              <Calendar className='mr-2 h-4 w-4' />
              Update Schedule
            </Button>
          )}
        </div>
      )
    },
  },
]
