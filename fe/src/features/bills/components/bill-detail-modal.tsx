import { Dialog, DialogContent, DialogHeader, DialogTitle, Badge } from '@/components/ui'
import { Receipt, Paperclip, Image as ImageIcon } from 'lucide-react'
import type { BillType } from '../model/bill-types'
import { BillQR } from './bill-qr'

interface BillDetailModalProps {
  bill: BillType | null
  open: boolean
  onClose: () => void
  onStatusChange?: (billId: number, newStatus: 'Unpaid' | 'Paid') => void
}

export const BillDetailModal = ({ bill, open, onClose }: BillDetailModalProps) => {
  if (!bill) return null

  // const toggleStatus = () => {
  //   if (bill && onStatusChange) {
  //     const newStatus = bill.statusBill === 'Paid' ? 'Unpaid' : 'Paid'
  //     onStatusChange(bill.id, newStatus)
  //     toast.success(`Bill marked as ${newStatus}`)
  //   }
  // }

  const getStatusColor = (status: string) => {
    const colors = {
      Unpaid: 'bg-red-100 text-red-800 dark:bg-red-400 dark:text-red-700',
      Paid: 'bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-200',
    } as const
    return colors[status as keyof typeof colors] || colors.Unpaid
  }

  // const formatDate = (dateString: string) => {
  //   try {
  //     if (!dateString) return 'N/A'
  //     return new Date(dateString).toLocaleDateString('en-US', {
  //       year: 'numeric',
  //       month: 'short',
  //       day: 'numeric',
  //       hour: '2-digit',
  //       minute: '2-digit',
  //     })
  //   } catch {
  //     return 'N/A'
  //   }
  // }

  const formatDateShort = (dateString: string) => {
    try {
      if (!dateString) return 'N/A'
      return new Date(dateString).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
      })
    } catch {
      return 'N/A'
    }
  }

  const getAssetName = () => {
    return bill.assets?.assetName || 'Unknown Asset'
  }

  const getCategoryName = () => {
    return bill.assets?.category?.categoryName || 'No Category'
  }

  const getAssetCost = () => {
    if (bill.assets?.cost !== undefined && bill.assets.cost !== null) {
      return bill.assets.cost
    }
    return bill.amount || 0
  }

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen) {
      onClose()
    }
  }

  // const getCreatorInitials = () => {
  //   const name = bill.creator?.fullName || 'Unknown User'
  //   return name
  //     .split(' ')
  //     .map((n) => n[0])
  //     .join('')
  //     .toUpperCase()
  //     .slice(0, 2)
  // }

  // const isWarrantyExpired = (warrantyDate?: string) => {
  //   if (!warrantyDate) return false
  //   return new Date(warrantyDate) < new Date()
  // }

  // const getLastUpdated = () => {
  //   if (!bill.updateAt || bill.updateAt === bill.createAt) {
  //     return bill.createAt
  //   }
  //   return bill.updateAt
  // }

  const hasFileAttachment = () => {
    return bill.fileAttachment && bill.fileAttachment.trim() !== '' && bill.fileAttachment !== 'null'
  }

  const hasImageUpload = () => {
    return bill.imageUpload && bill.imageUpload.trim() !== '' && bill.imageUpload !== 'null'
  }

  const openFile = (url: string, type: 'file' | 'image') => {
    if (url && url !== 'null') {
      try {
        window.open(url, '_blank')
      } catch (error) {
        console.error(`Failed to open ${type}:`, error)
      }
    }
  }

  // const InfoRow = ({
  //   label,
  //   value,
  //   valueClassName = '',
  //   badge = false,
  //   badgeClassName = '',
  // }: {
  //   label: string
  //   value: string | number
  //   valueClassName?: string
  //   badge?: boolean
  //   badgeClassName?: string
  // }) => (
  //   <div className='flex items-start gap-3 border-gray-100 py-3 last:border-b-0 dark:border-gray-700'>
  //     <div className='text-sm text-gray-700 dark:text-gray-300'>
  //       <span className='font-medium'>{label}: </span>
  //       {badge ? (
  //         <Badge
  //           className={`${badgeClassName} ml-1`}
  //           variant='outline'
  //         >
  //           {value}
  //         </Badge>
  //       ) : (
  //         <span className={`text-gray-900 dark:text-gray-100 ${valueClassName}`}>{value}</span>
  //       )}
  //     </div>
  //   </div>
  // )

  const sellerInfo = {
    companyName: 'SUNRISE SOFTWARE SOLUTIONS CORPORATION',
    taxCode: '0305089644',
    address: '307/12 Nguyen Van Troi, Ward 1, Tan Binh District, HCMC, Viet Nam',
    phoneNumber: '028-35471411',
    accountNumber: '1982738238232',
  }

  const buyerInfo = {
    name: bill.creator?.fullName || 'John Doe',
    taxCode: '',
    address: '533/8 Nguyen Tri Phuong, Ward 8, District 10, Ho Chi Minh City',
    phoneNumber: '032-35471411',
    accountNumber: '4568239472356',
    paymentMethod: 'Bank Transfer',
  }

  const items = [
    {
      stt: 1,
      assetName: getAssetName(),
      category: getCategoryName(),
      quantity: 1,
      unitPrice: getAssetCost(),
      amount: getAssetCost(),
    },
  ]

  const totalAmount = items.reduce((sum, item) => sum + item.amount, 0)

  return (
    <Dialog
      open={open}
      onOpenChange={handleOpenChange}
      modal={true}
    >
      <DialogContent className='!h-[90vh] !w-[50vw] !max-w-[1200px] overflow-y-auto bg-white dark:bg-gray-900'>
        <DialogHeader className='border-b pb-4 dark:border-gray-700'>
          <div className='flex items-center justify-between'>
            <DialogTitle className='flex items-center gap-3 text-xl font-bold text-gray-900 dark:text-gray-100'>
              <Receipt className='h-5 w-5' />
              SALES INVOICE - {bill.billNumber}
            </DialogTitle>
          </div>
        </DialogHeader>

        <div className='rounded-lg border border-gray-200 bg-white p-6 dark:border-gray-700 dark:bg-gray-800'>
          <div className='flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between'>
            <div className='flex-shrink-0'>
              <img
                src='https://www.s3corp.com.vn/images/S3CORP.svg?w=128&q=75'
                alt='Company Logo'
                className='h-14 w-auto object-contain'
              />
            </div>

            <div className='flex flex-1 flex-col items-center justify-center text-center'>
              <div className='bg-red-500 px-4 py-2 text-lg font-bold text-white'>SALES INVOICE</div>
              <p className='mt-1 text-sm text-gray-700 dark:text-gray-300'>
                Date: <span className='font-semibold'>{formatDateShort(bill.createAt)}</span>
              </p>
            </div>

            <div className='flex items-start gap-4'>
              <div className='flex flex-col items-end justify-between gap-1'>
                <p className='text-sm text-gray-600 dark:text-gray-100'>
                  Form No.: <span className='font-semibold'>2C24TTU</span>
                </p>
                <p className='text-sm text-gray-600 dark:text-gray-100'>
                  Invoice No.: <span className='font-semibold'>{bill.billNumber}</span>
                </p>
              </div>

              <div>
                <BillQR bill={bill} />
              </div>
            </div>
          </div>
        </div>

        <div className='rounded-lg border border-gray-200 bg-white p-4 dark:border-gray-700 dark:bg-gray-800'>
          <h3 className='mb-3 text-sm font-semibold text-gray-900 dark:text-gray-100'>
            Seller: <span className='font-normal'>{sellerInfo.companyName}</span>
          </h3>
          <div className='grid grid-cols-1 gap-2 text-sm'>
            <h3 className='font-semibold text-gray-900 dark:text-gray-100'>
              Tax Code: <span className='font-normal'>{sellerInfo.taxCode}</span>
            </h3>
            <h3 className='font-semibold text-gray-900 dark:text-gray-100'>
              Address: <span className='font-normal'>{sellerInfo.address}</span>
            </h3>
            <h3 className='font-semibold text-gray-900 dark:text-gray-100'>
              Phone: <span className='font-normal'>{sellerInfo.phoneNumber}</span>
            </h3>
            <h3 className='border-b pb-2 font-semibold text-gray-900 dark:text-gray-100'>
              Account No.: <span className='font-normal'>{sellerInfo.accountNumber}</span>
            </h3>
            <h3 className='mt-0.5 text-sm font-semibold text-gray-900 dark:text-gray-100'>
              Buyer: <span className='font-normal'>{buyerInfo.name}</span>
            </h3>
            <h3 className='font-semibold text-gray-900 dark:text-gray-100'>
              Tax Code: <span className='font-normal'>{buyerInfo.taxCode || '_________________'}</span>
            </h3>
            <h3 className='font-semibold text-gray-900 dark:text-gray-100'>
              Address: <span className='font-normal'>{buyerInfo.address}</span>
            </h3>
            <h3 className='font-semibold text-gray-900 dark:text-gray-100'>
              Phone: <span className='font-normal'>{buyerInfo.phoneNumber}</span>
            </h3>
            <h3 className='font-semibold text-gray-900 dark:text-gray-100'>
              Account No.: <span className='font-normal'>{buyerInfo.accountNumber}</span>
            </h3>
            <h3 className='font-semibold text-gray-900 dark:text-gray-100'>
              Payment Method: <span className='font-normal'>{buyerInfo.paymentMethod}</span>
            </h3>
          </div>
        </div>

        <div className='rounded-lg border border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800'>
          <div className='overflow-x-auto'>
            <table className='w-full text-sm'>
              <thead className='bg-gray-50 dark:bg-gray-700'>
                <tr className='border-b dark:border-gray-600'>
                  <th className='border-r px-3 py-2 text-center text-xs font-semibold text-gray-900 dark:border-gray-600 dark:text-gray-300'>
                    No.
                  </th>
                  <th className='border-r px-3 py-2 text-center text-xs font-semibold text-gray-900 dark:border-gray-600 dark:text-gray-300'>
                    Asset Name
                  </th>
                  <th className='border-r px-3 py-2 text-center text-xs font-semibold text-gray-900 dark:border-gray-600 dark:text-gray-300'>
                    Category
                  </th>
                  <th className='border-r px-3 py-2 text-center text-xs font-semibold text-gray-900 dark:border-gray-600 dark:text-gray-300'>
                    Quantity
                  </th>
                  <th className='border-r px-3 py-2 text-center text-xs font-semibold text-gray-900 dark:border-gray-600 dark:text-gray-300'>
                    Unit Price
                  </th>
                  <th className='px-3 py-2 text-center text-xs font-semibold text-gray-900 dark:text-gray-300'>
                    Amount
                  </th>
                </tr>
              </thead>
              <tbody className='bg-white dark:bg-gray-800'>
                {items.map((item, index) => (
                  <tr
                    key={index}
                    className='border-b dark:border-gray-600'
                  >
                    <td className='border-r px-3 py-2 text-center dark:border-gray-600'>{item.stt}</td>
                    <td className='border-r px-3 py-2 text-center dark:border-gray-600'>{item.assetName}</td>
                    <td className='border-r px-3 py-2 text-center dark:border-gray-600'>{item.category}</td>
                    <td className='border-r px-3 py-2 text-center dark:border-gray-600'>{item.quantity}</td>
                    <td className='border-r px-3 py-2 text-center dark:border-gray-600'>
                      ${item.unitPrice.toLocaleString()}
                    </td>
                    <td className='px-3 py-2 text-center'>${item.amount.toLocaleString()}</td>
                  </tr>
                ))}
                {[...Array(3)].map((_, index) => (
                  <tr
                    key={`empty-${index}`}
                    className='border-b dark:border-gray-600'
                  >
                    <td className='border-r px-3 py-2 text-center dark:border-gray-600'>&nbsp;</td>
                    <td className='border-r px-3 py-2 dark:border-gray-600'>&nbsp;</td>
                    <td className='border-r px-3 py-2 text-center dark:border-gray-600'>&nbsp;</td>
                    <td className='border-r px-3 py-2 text-center dark:border-gray-600'>&nbsp;</td>
                    <td className='border-r px-3 py-2 text-right dark:border-gray-600'>&nbsp;</td>
                    <td className='px-3 py-2 text-right'>&nbsp;</td>
                  </tr>
                ))}
              </tbody>
              <tfoot className='bg-gray-50 dark:bg-gray-700'>
                <tr className='border-b dark:border-gray-600'>
                  <td
                    colSpan={2}
                    className='px-3 py-2 text-left text-sm font-semibold text-gray-900 dark:text-gray-300'
                  >
                    <span className='px-3 py-2 text-right dark:border-gray-600'>Status:</span>
                    <Badge className={getStatusColor(bill.statusBill || 'Unpaid')}>
                      {bill.statusBill === 'Paid' ? 'Paid' : 'Unpaid'}
                    </Badge>
                  </td>
                  <td
                    colSpan={3}
                    className='px-3 py-2 text-right font-semibold dark:border-gray-600'
                  >
                    Total Payment:
                  </td>
                  <td className='px-3 py-2 text-center font-bold text-green-600'>${totalAmount.toLocaleString()}</td>
                </tr>
              </tfoot>
            </table>
          </div>
        </div>
        <div className='px-2 py-0 text-sm text-gray-700 italic dark:text-gray-300'>
          <p className='font-semibold'>Note:</p>
          <ul className='list-disc space-y-1 pl-5'>
            <li>This invoice is valid only when signed and stamped by the seller.</li>
            <li>Please retain this invoice for warranty and accounting purposes.</li>
            <li>All prices are inclusive of applicable taxes (if any).</li>
            <li>Goods/services sold are non-refundable unless otherwise stated.</li>
            <li>Contact our support team if any discrepancies are found.</li>
          </ul>
        </div>

        <div className='rounded-lg border border-gray-200 bg-white p-3 dark:border-gray-700 dark:bg-gray-800'>
          <div className='grid grid-cols-2 gap-10 pt-6 text-sm text-gray-700 dark:text-gray-300'>
            <div className='text-center'>
              <p className='mb-1 font-semibold text-gray-900 dark:text-gray-100'>Buyer</p>
              <p className='mb-8 text-xs'>(Signature and full name)</p>
              <div className='border-gray-400 pb-2'>
                <p className='font-semibold'>_________________</p>
              </div>
            </div>
            <div className='text-center'>
              <p className='mb-1 font-semibold text-gray-900 dark:text-gray-100'>Seller</p>
              <p className='mb-8 text-xs'>(Signature and full name)</p>
              <div className='border-gray-400 pb-2'>
                <p className='font-semibold'>{bill.creator?.fullName}</p>
              </div>
            </div>
          </div>
        </div>

        {(hasFileAttachment() || hasImageUpload()) && (
          <div className='rounded-lg border border-gray-200 bg-white p-6 dark:border-gray-700 dark:bg-gray-800'>
            <h3 className='mb-4 flex items-center gap-2 text-lg font-semibold text-gray-900 dark:text-gray-100'>
              <Paperclip className='h-5 w-5' />
              Attachments
            </h3>
            <div className='grid grid-cols-1 gap-4 lg:grid-cols-2'>
              {hasFileAttachment() && (
                <div className='flex items-center gap-3'>
                  <Paperclip className='h-4 w-4 text-gray-500' />
                  <button
                    onClick={() => openFile(bill.fileAttachment!, 'file')}
                    className='text-blue-600 hover:text-blue-800 hover:underline dark:text-blue-400'
                  >
                    View Document
                  </button>
                </div>
              )}
              {hasImageUpload() && (
                <div className='flex items-center gap-3'>
                  <ImageIcon className='h-4 w-4 text-gray-500' />
                  <button
                    onClick={() => openFile(bill.imageUpload!, 'image')}
                    className='text-blue-600 hover:text-blue-800 hover:underline dark:text-blue-400'
                  >
                    View Image
                  </button>
                </div>
              )}
            </div>
          </div>
        )}
      </DialogContent>
    </Dialog>
  )
}

export default BillDetailModal
