import { useState } from 'react'
import { Button } from '@/components/ui'
import { Printer, Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import type { BillType } from '../model/bill-types'
import { BillPrintLayout } from './bill-print-layout'

interface BillPrintButtonProps {
  bill: BillType
  onStatusUpdated?: (newStatus: string) => void
  variant?: 'default' | 'outline' | 'ghost'
  size?: 'default' | 'sm' | 'lg'
}

export const BillPrintButton = ({
  bill,
  onStatusUpdated,
  variant = 'default',
  size = 'default',
}: BillPrintButtonProps) => {
  const [isPrinting, setIsPrinting] = useState(false)

  const handlePrint = async () => {
    setIsPrinting(true)

    try {
      const printWindow = window.open('', '_blank', 'width=1200,height=800')

      if (!printWindow) {
        toast.error('Unable to open print window. Please allow popups for this site.')
        setIsPrinting(false)
        return
      }

      const printDocument = `
        <!DOCTYPE html>
        <html lang="en">
          <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Bill ${bill.billNumber}</title>
            <style>
              * {
                margin: 0;
                padding: 0;
                box-sizing: border-box;
              }
              
              body {
                font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
                line-height: 1.6;
                color: #000;
                background: white;
              }
              
              .bill-print-container {
                padding: 20px;
                max-width: 1200px;
                margin: 0 auto;
              }
              
              @media print {
                body {
                  background: white !important;
                  -webkit-print-color-adjust: exact;
                  print-color-adjust: exact;
                }
                
                .bill-print-container {
                  padding: 0;
                  max-width: none;
                  margin: 0;
                }
                
                .no-print {
                  display: none !important;
                }
                
                .rounded-lg {
                  border-radius: 0;
                }
                
                .border-gray-200 {
                  border-color: #000 !important;
                }
                
                .text-gray-500, .text-gray-600, .text-gray-700 {
                  color: #000 !important;
                }
                
                .bg-gray-50, .bg-gray-100 {
                  background-color: #f5f5f5 !important;
                }
                
                img {
                  max-width: 100% !important;
                  height: auto !important;
                }
                
                .grid {
                  display: block !important;
                }
                
                .lg\\:col-span-3, .lg\\:col-span-1 {
                  width: 100% !important;
                  margin-bottom: 20px;
                }
                
                .space-y-6 > * + * {
                  margin-top: 20px !important;
                }
                
                .break-inside-avoid {
                  break-inside: avoid;
                }
              }
              
              .rounded-lg {
                border-radius: 8px;
              }
              
              .border {
                border-width: 1px;
              }
              
              .border-gray-200 {
                border-color: #e5e7eb;
              }
              
              .bg-white {
                background-color: white;
              }
              
              .p-6 {
                padding: 24px;
              }
              
              .mb-4 {
                margin-bottom: 16px;
              }
              
              .flex {
                display: flex;
              }
              
              .items-center {
                align-items: center;
              }
              
              .gap-2 {
                gap: 8px;
              }
              
              .gap-3 {
                gap: 12px;
              }
              
              .gap-4 {
                gap: 16px;
              }
              
              .text-lg {
                font-size: 18px;
              }
              
              .font-semibold {
                font-weight: 600;
              }
              
              .font-bold {
                font-weight: 700;
              }
              
              .font-medium {
                font-weight: 500;
              }
              
              .text-gray-900 {
                color: #111827;
              }
              
              .text-gray-700 {
                color: #374151;
              }
              
              .text-gray-600 {
                color: #4b5563;
              }
              
              .text-gray-500 {
                color: #6b7280;
              }
              
              .text-sm {
                font-size: 14px;
              }
              
              .text-xs {
                font-size: 12px;
              }
              
              .space-y-0 > * + * {
                margin-top: 0;
              }
              
              .space-y-6 > * + * {
                margin-top: 24px;
              }
              
              .py-3 {
                padding-top: 12px;
                padding-bottom: 12px;
              }
              
              .border-b {
                border-bottom-width: 1px;
              }
              
              .border-gray-100 {
                border-color: #f3f4f6;
              }
              
              .grid {
                display: grid;
              }
              
              .grid-cols-1 {
                grid-template-columns: repeat(1, minmax(0, 1fr));
              }
              
              .grid-cols-2 {
                grid-template-columns: repeat(2, minmax(0, 1fr));
              }
              
              .gap-6 {
                gap: 24px;
              }
              
              .gap-10 {
                gap: 40px;
              }
              
              .lg\\:grid-cols-4 {
                grid-template-columns: repeat(4, minmax(0, 1fr));
              }
              
              .lg\\:col-span-3 {
                grid-column: span 3 / span 3;
              }
              
              .lg\\:col-span-1 {
                grid-column: span 1 / span 1;
              }
              
              .text-center {
                text-align: center;
              }
              
              .text-right {
                text-align: right;
              }
              
              .mt-3 {
                margin-top: 12px;
              }
              
              .mt-6 {
                margin-top: 24px;
              }
              
              .mb-12 {
                margin-bottom: 48px;
              }
              
              .pt-6 {
                padding-top: 24px;
              }
              
              .pb-1 {
                padding-bottom: 4px;
              }
              
              .border-gray-400 {
                border-color: #9ca3af;
              }
              
              .h-12 {
                height: 48px;
              }
              
              .w-auto {
                width: auto;
              }
              
              .object-contain {
                object-fit: contain;
              }
              
              .h-4 {
                height: 16px;
              }
              
              .w-4 {
                width: 16px;
              }
              
              .h-5 {
                height: 20px;
              }
              
              .w-5 {
                width: 20px;
              }
              
              .mt-0\\.5 {
                margin-top: 2px;
              }
              
              .flex-1 {
                flex: 1 1 0%;
              }
              
              .break-words {
                overflow-wrap: break-word;
              }
              
              .break-all {
                word-break: break-all;
              }
              
              .whitespace-pre-wrap {
                white-space: pre-wrap;
              }
              
              .leading-relaxed {
                line-height: 1.625;
              }
              
              .font-mono {
                font-family: ui-monospace, SFMono-Regular, "SF Mono", Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
              }
              
              .text-green-600 {
                color: #059669;
              }
              
              .text-red-600 {
                color: #dc2626;
              }
              
              .text-yellow-500 {
                color: #eab308;
              }
              
              .text-blue-600 {
                color: #2563eb;
              }
              
              .bg-gray-100 {
                background-color: #f3f4f6;
              }
              
              .bg-yellow-100 {
                background-color: #fef3c7;
              }
              
              .bg-blue-100 {
                background-color: #dbeafe;
              }
              
              .bg-green-100 {
                background-color: #dcfce7;
              }
              
              .bg-red-100 {
                background-color: #fee2e2;
              }
              
              .text-gray-800 {
                color: #1f2937;
              }
              
              .text-yellow-800 {
                color: #92400e;
              }
              
              .text-blue-800 {
                color: #1e40af;
              }
              
              .text-green-800 {
                color: #166534;
              }
              
              .text-red-800 {
                color: #991b1b;
              }
              
              .inline-block {
                display: inline-block;
              }
              
              .px-2 {
                padding-left: 8px;
                padding-right: 8px;
              }
              
              .py-1 {
                padding-top: 4px;
                padding-bottom: 4px;
              }
              
              .rounded {
                border-radius: 4px;
              }
              
              .ml-1 {
                margin-left: 4px;
              }
              
              .badge {
                display: inline-block;
                padding: 4px 8px;
                border-radius: 4px;
                font-size: 12px;
                font-weight: 500;
                border: 1px solid currentColor;
              }
              
              .avatar {
                display: inline-flex;
                align-items: center;
                justify-content: center;
                vertical-align: middle;
                overflow: hidden;
                border-radius: 50%;
                background-color: #f3f4f6;
                color: #111827;
                font-weight: 500;
              }
              
              .avatar img {
                width: 100%;
                height: 100%;
                object-fit: cover;
              }
              
              .qr-container {
                max-width: 200px;
                margin: 0 auto;
              }
              
              .qr-container img {
                width: 100%;
                height: auto;
                border-radius: 8px;
              }
              
              @page {
                margin: 20mm;
                size: A4;
              }
            </style>
          </head>
          <body>
            <div id="print-content"></div>
            <script>
              // Print when page loads
              window.onload = function() {
                setTimeout(function() {
                  window.print();
                }, 500);
              };
              
              // Close window after printing
              window.onafterprint = function() {
                window.close();
              };
            </script>
          </body>
        </html>
      `

      printWindow.document.write(printDocument)
      printWindow.document.close()

      const tempDiv = document.createElement('div')
      document.body.appendChild(tempDiv)

      const { createRoot } = await import('react-dom/client')
      const root = createRoot(tempDiv)

      await new Promise<void>((resolve) => {
        root.render(<BillPrintLayout bill={bill} />)
        setTimeout(() => {
          const printContent = tempDiv.innerHTML

          const printContentElement = printWindow.document.getElementById('print-content')
          if (printContentElement) {
            printContentElement.innerHTML = printContent
          }

          root.unmount()
          document.body.removeChild(tempDiv)

          resolve()
        }, 100)
      })

      if (onStatusUpdated && bill.status !== 'Paid') {
        onStatusUpdated('Paid')
        toast.success('Bill status updated to Paid')
      }

      toast.success('Print window opened successfully')
    } catch (error) {
      console.error('Print error:', error)
      toast.error('Failed to print bill')
    } finally {
      setIsPrinting(false)
    }
  }

  return (
    <Button
      onClick={handlePrint}
      disabled={isPrinting}
      variant={variant}
      size={size}
      className='flex items-center gap-2'
    >
      {isPrinting ? <Loader2 className='h-4 w-4 animate-spin' /> : <Printer className='h-4 w-4' />}
      {isPrinting ? 'Printing...' : 'Print Bill'}
    </Button>
  )
}
