import { Eye, Printer, Calendar } from 'lucide-react'
import { Button, Badge } from '@/components/ui'
import { DataTable } from '@/components/ui/data-table-component'
import type { BillType } from '../model/bill-types'
import { toast } from 'sonner'

interface BillsTableProps {
  bills: BillType[]
  onViewBill: (bill: BillType) => void
  isLoading: boolean
}

export const BillsTable = ({ bills, onViewBill, isLoading }: BillsTableProps) => {
  const handlePrintBill = (bill: BillType) => {
    try {
      const printContent = `
        <!DOCTYPE html>
        <html>
          <head>
            <title>Bill ${bill.billNumber}</title>
            <style>
              body { 
                font-family: Arial, sans-serif; 
                padding: 40px; 
                max-width: 800px; 
                margin: 0 auto;
                line-height: 1.6;
              }
              .header { 
                text-align: center; 
                border-bottom: 2px solid #333; 
                padding-bottom: 20px; 
                margin-bottom: 30px; 
              }
              .bill-info { 
                display: grid; 
                grid-template-columns: 1fr 1fr; 
                gap: 20px; 
                margin-bottom: 30px; 
              }
              .info-row { margin: 8px 0; }
              .label { font-weight: bold; color: #333; }
              .amount { 
                font-size: 28px; 
                font-weight: bold; 
                color: #007bff; 
                text-align: center;
                padding: 20px;
                background: #f8f9fa;
                border-radius: 8px;
                margin: 20px 0;
              }
              .description-box {
                border: 1px solid #ddd;
                padding: 15px;
                background: #f9f9f9;
                border-radius: 5px;
                margin: 20px 0;
              }
              @media print {
                body { margin: 0; }
                .no-print { display: none; }
              }
            </style>
          </head>
          <body>
            <div class="header">
              <h1>BILL INVOICE</h1>
              <p>Generated on ${new Date().toLocaleDateString()}</p>
            </div>
            
            <div class="bill-info">
              <div>
                <div class="info-row">
                  <span class="label">Bill Number:</span> ${bill.billNumber || 'N/A'}
                </div>
                <div class="info-row">
                  <span class="label">Asset Name:</span> ${getAssetName(bill)}
                </div>
                <div class="info-row">
                  <span class="label">Category:</span> ${getCategoryName(bill)}
                </div>
              </div>
              <div>
                <div class="info-row">
                  <span class="label">Status:</span> ${bill.status || 'Unpaid'}
                </div>
                <div class="info-row">
                  <span class="label">Created Date:</span> ${formatDateForPrint(bill.createAt)}
                </div>
                <div class="info-row">
                  <span class="label">Updated Date:</span> ${formatDateForPrint(bill.updateAt)}
                </div>
              </div>
            </div>

            <div class="description-box">
              <div class="label">Description:</div>
              <p>${bill.description || 'No description provided'}</p>
            </div>

            <div class="amount">
              Asset Cost: $${getAssetCost(bill).toLocaleString()}
            </div>
          </body>
        </html>
      `

      const newWindow = window.open('', '_blank')
      if (newWindow) {
        newWindow.document.write(printContent)
        newWindow.document.close()
        setTimeout(() => {
          newWindow.print()
        }, 100)
        toast.success('Print window opened')
      } else {
        toast.error('Please allow popups for printing')
      }
    } catch (error) {
      console.error('Print error:', error)
      toast.error('Print failed')
    }
  }

  const formatDateForPrint = (dateString: string) => {
    try {
      if (!dateString) return new Date().toLocaleDateString()
      return new Date(dateString).toLocaleDateString()
    } catch {
      return new Date().toLocaleDateString()
    }
  }

  const getStatusBadge = (status: string) => {
    const statusConfig = {
      Unpaid: { variant: 'secondary', color: 'bg-gray-100 text-gray-800' },

      Paid: { variant: 'success', color: 'bg-green-100 text-green-800' },
    } as const

    const config = statusConfig[status as keyof typeof statusConfig] || statusConfig.Unpaid

    return <Badge className={config.color}>{status || 'Unpaid'}</Badge>
  }

  const formatDate = (dateString: string) => {
    try {
      if (!dateString) {
        return new Date().toLocaleDateString()
      }
      return new Date(dateString).toLocaleDateString()
    } catch {
      return new Date().toLocaleDateString()
    }
  }

  const formatCurrency = (amount: number) => {
    try {
      return `$${(amount || 0).toLocaleString()}`
    } catch {
      return '$0'
    }
  }

  const getAssetName = (bill: BillType) => {
    if (bill.assets?.assetName) {
      return bill.assets.assetName
    }
    return 'Unknown Asset'
  }

  const getCategoryName = (bill: BillType) => {
    console.log('Debug bill:', bill)
    if (bill.assets?.category?.categoryName) {
      return bill.assets.category.categoryName
    }
    return 'No Category'
  }

  const getAssetCost = (bill: BillType) => {
    if (bill.assets?.cost !== undefined && bill.assets.cost !== null) {
      return bill.assets.cost
    }
    return bill.amount || 0
  }

  const columns = [
    {
      accessorKey: 'billNumber',
      header: 'Bill Number',
      cell: ({ row }: any) => <div className='font-medium'>{row.original.billNumber || `BILL-${row.original.id}`}</div>,
    },
    {
      accessorKey: 'assets.assetName',
      header: 'Asset Name',
      cell: ({ row }: any) => (
        <div
          className='max-w-[150px] truncate font-medium'
          title={getAssetName(row.original)}
        >
          {getAssetName(row.original)}
        </div>
      ),
    },
    {
      accessorKey: 'assets.category.categoryName',
      header: 'Category',
      cell: ({ row }: any) => <Badge variant='outline'>{getCategoryName(row.original)}</Badge>,
    },
    {
      accessorKey: 'description',
      header: 'Description',
      cell: ({ row }: any) => (
        <div
          className='max-w-[200px] truncate'
          title={row.original.description}
        >
          {row.original.description || 'No description'}
        </div>
      ),
    },
    {
      accessorKey: 'assets.cost',
      header: 'Cost',
      cell: ({ row }: any) => (
        <div className='font-semibold text-green-600'>{formatCurrency(getAssetCost(row.original))}</div>
      ),
    },
    {
      accessorKey: 'status',
      header: 'Status',
      cell: ({ row }: any) => getStatusBadge(row.original.status),
    },
    {
      accessorKey: 'createdAt',
      header: 'Created Date',
      cell: ({ row }: any) => (
        <div className='text-muted-foreground flex items-center gap-1 text-sm'>
          <Calendar className='h-3 w-3' />
          {formatDate(row.original.createdAt)}
        </div>
      ),
    },
    {
      id: 'actions',
      header: 'Actions',
      cell: ({ row }: any) => (
        <div className='flex items-center gap-2'>
          <Button
            variant='ghost'
            size='sm'
            onClick={() => onViewBill(row.original)}
            className='h-8 w-8 p-0'
            title='View details'
          >
            <Eye className='h-4 w-4' />
          </Button>
          <Button
            variant='ghost'
            size='sm'
            onClick={() => handlePrintBill(row.original)}
            className='h-8 w-8 p-0'
            title='Print bill'
          >
            <Printer className='h-4 w-4' />
          </Button>
        </div>
      ),
    },
  ]

  return (
    <div className='rounded-md border'>
      <DataTable
        columns={columns}
        data={bills || []}
        isLoading={isLoading}
        emptyMessage='No bills found. Create your first bill to get started.'
      />
    </div>
  )
}
