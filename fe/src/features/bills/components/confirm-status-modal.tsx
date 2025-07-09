import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { AlertCircle } from 'lucide-react'

interface ConfirmStatusModalProps {
  open: boolean
  onClose: () => void
  onConfirm: () => void
  billNumber: string
}

export const ConfirmStatusModal = ({ open, onClose, onConfirm, billNumber }: ConfirmStatusModalProps) => {
  return (
    <Dialog
      open={open}
      onOpenChange={onClose}
    >
      <DialogContent>
        <DialogHeader>
          <DialogTitle>
            <div className='flex items-center gap-2'>
              <AlertCircle className='text-600 h-5 w-5' />
              <span>Confirm Status Change</span>
            </div>
          </DialogTitle>
          <DialogDescription>
            Are you sure you want to mark bill #{billNumber} as paid? This action cannot be undone.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button
            variant='outline'
            onClick={onClose}
          >
            Cancel
          </Button>
          <Button onClick={onConfirm}>Confirm</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
