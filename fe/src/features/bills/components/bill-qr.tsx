import type { BillType } from '../model/bill-types'
import { QRCodeSVG } from 'qrcode.react'

interface BillQRProps {
  bill: BillType
}

export const BillQR = ({ bill }: BillQRProps) => {
  const billUrl = `${window.location.origin}/bills/${bill.billNumber}`

  if (bill.qrUrl) {
    return (
      <img
        src={bill.qrUrl}
        alt={`QR Code for ${bill.billNumber}`}
        className='max-h-[200px] w-full rounded-md object-contain'
      />
    )
  }

  return (
    <div className='flex h-[200px] w-full flex-col items-center justify-center rounded-md border border-dashed bg-gray-50'>
      <QRCodeSVG
        value={billUrl}
        size={80}
      />
      <p className='mt-2 px-2 text-center text-xs break-all text-gray-400'>{billUrl}</p>
    </div>
  )
}
