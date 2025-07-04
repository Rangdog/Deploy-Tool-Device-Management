import { useState } from 'react'
import { createRoot } from 'react-dom/client'
import { Button } from '@/components/ui'
import { Printer } from 'lucide-react'
import { toast } from 'sonner'
import { BillPrintLayout } from './bill-print-layout'
import type { BillType } from '../model/bill-types'

interface BillPrintButtonProps {
  bill: BillType
  variant?: 'default' | 'outline' | 'ghost'
  size?: 'sm' | 'default' | 'lg'
  className?: string
}

export const BillPrintButton = ({ bill, variant = 'outline', size = 'sm', className = '' }: BillPrintButtonProps) => {
  const [isPrinting, setIsPrinting] = useState(false)

  const handlePrint = async () => {
    setIsPrinting(true)

    try {
      const printWindow = window.open('', '_blank', 'width=1200,height=800')

      if (!printWindow) {
        toast.error('Please allow popups for printing')
        setIsPrinting(false)
        return
      }

      printWindow.document.write(`
        <!DOCTYPE html>
        <html lang="en">
          <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Bill ${bill.billNumber}</title>
            <script src="https://unpkg.com/react@18/umd/react.development.js"></script>
            <script src="https://unpkg.com/react-dom@18/umd/react-dom.development.js"></script>
            <script src="https://cdn.tailwindcss.com"></script>
            <style>
              @media print {
                .no-print { display: none !important; }
                body { margin: 0; padding: 0; }
                @page { margin: 1cm; size: A4; }
              }
            </style>
          </head>
          <body>
            <div id="print-root"></div>
          </body>
        </html>
      `)

      printWindow.document.close()

      // Wait for scripts to load
      setTimeout(() => {
        const container = printWindow.document.getElementById('print-root')
        if (container) {
          const root = createRoot(container)
          root.render(<BillPrintLayout bill={bill} />)

          setTimeout(() => {
            printWindow.print()
            toast.success('Print completed')
          }, 1000)
        }
      }, 500)
    } catch (error) {
      console.error('Print error:', error)
      toast.error('Failed to print bill')
    }

    setIsPrinting(false)
  }

  return (
    <Button
      variant={variant}
      size={size}
      onClick={handlePrint}
      disabled={isPrinting}
      className={className}
    >
      <Printer className='mr-1 h-4 w-4' />
      {isPrinting ? 'Printing...' : 'Print Bill'}
    </Button>
  )
}
