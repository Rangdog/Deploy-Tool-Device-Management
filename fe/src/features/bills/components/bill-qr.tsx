import { QrCode } from 'lucide-react'
import type { BillType } from '../model/bill-types'

interface BillQRProps {
  bill: BillType
}

export const BillQR = ({ bill }: BillQRProps) => {
  if (bill.qrUrl) {
    return (
      <img
        src={bill.qrUrl}
        alt={`QR Code for ${bill.billNumber}`}
        className='max-h-[200px] w-full rounded-md object-contain'
      />
    )
  }

  const billUrl = `${window.location.origin}/bills/${bill.id}`

  return (
    <div className='flex h-[200px] w-full flex-col items-center justify-center rounded-md border border-dashed bg-gray-50'>
      <QrCode className='text-muted-foreground mb-2 h-10 w-10' />
      <p className='text-muted-foreground px-4 text-center text-sm'>QR Code will be generated</p>
      <p className='mt-1 px-2 text-center text-xs break-all text-gray-400'>{billUrl}</p>
    </div>
  )
}
