import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  Badge,
  Avatar,
  AvatarFallback,
  AvatarImage,
} from '@/components/ui'
import {
  Receipt,
  Calendar,
  DollarSign,
  FileText,
  Building2,
  Tag,
  User,
  Mail,
  Hash,
  Shield,
  Paperclip,
  Image as ImageIcon,
  QrCode,
  Package,
  Clock,
  AlertCircle,
  ExternalLink,
} from 'lucide-react'
import type { BillType } from '../model/bill-types'
import { BillQR } from './bill-qr'
interface BillDetailModalProps {
  bill: BillType | null
  open: boolean
  onClose: () => void
}

export const BillDetailModal = ({ bill, open, onClose }: BillDetailModalProps) => {
  if (!bill) return null

  const getStatusColor = (status: string) => {
    const colors = {
      Unpaid: 'bg-gray-100 text-gray-800',

      Paid: 'bg-green-100 text-green-800',
    } as const
    return colors[status as keyof typeof colors] || colors.Unpaid
  }

  const formatDate = (dateString: string) => {
    try {
      if (!dateString) return 'N/A'
      return new Date(dateString).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
      })
    } catch {
      return 'N/A'
    }
  }

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

  const getCreatorInitials = () => {
    const name = bill.creator?.fullName || 'Unknown User'
    return name
      .split(' ')
      .map((n) => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2)
  }

  const isWarrantyExpired = (warrantyDate?: string) => {
    if (!warrantyDate) return false
    return new Date(warrantyDate) < new Date()
  }

  const getLastUpdated = () => {
    if (!bill.updateAt || bill.updateAt === bill.createAt) {
      return bill.createAt
    }
    return bill.updateAt
  }

  const InfoRow = ({
    icon: Icon,
    label,
    value,
    valueClassName = '',
    badge = false,
    badgeClassName = '',
  }: {
    icon: any
    label: string
    value: string | number
    valueClassName?: string
    badge?: boolean
    badgeClassName?: string
  }) => (
    <div className='flex items-start gap-3 border-b border-gray-100 py-3 last:border-b-0'>
      <Icon className='mt-0.5 h-4 w-4 text-gray-500' />
      <div className='text-sm text-gray-700'>
        <span className='font-medium'>{label}: </span>
        {badge ? (
          <Badge
            className={`${badgeClassName} ml-1`}
            variant='outline'
          >
            {value}
          </Badge>
        ) : (
          <span className={`text-gray-900 ${valueClassName}`}>{value}</span>
        )}
      </div>
    </div>
  )

  // const DoubleInfoRow = ({
  //   icon1: Icon1,
  //   label1,
  //   value1,
  //   valueClassName1 = '',
  //   icon2: Icon2,
  //   label2,
  //   value2,
  //   valueClassName2 = '',
  // }: {
  //   icon1: React.ElementType
  //   label1: string
  //   value1: string | number
  //   valueClassName1?: string
  //   icon2: React.ElementType
  //   label2: string
  //   value2: string | number
  //   valueClassName2?: string
  // }) => (
  //   <div className='grid grid-cols-2 gap-6 border-b border-gray-100 py-3 last:border-b-0'>
  //     <div className='flex items-start gap-2'>
  //       <Icon1 className='mt-0.5 h-4 w-4 text-gray-500' />
  //       <div className='text-sm text-gray-700'>
  //         <span className='font-medium'>{label1}: </span>
  //         <span className={`text-gray-900 ${valueClassName1}`}>{value1}</span>
  //       </div>
  //     </div>
  //     <div className='flex items-start gap-2'>
  //       <Icon2 className='mt-0.5 h-4 w-4 text-gray-500' />
  //       <div className='text-sm text-gray-700'>
  //         <span className='font-medium'>{label2}: </span>
  //         <span className={`text-gray-900 ${valueClassName2}`}>{value2}</span>
  //       </div>
  //     </div>
  //   </div>
  // )

  return (
    <Dialog
      open={open}
      onOpenChange={handleOpenChange}
    >
      <DialogContent className='!h-[90vh] !w-[50vw] !max-w-[1200px] overflow-y-auto bg-white'>
        <DialogHeader className='border-b pb-4'>
          <DialogTitle className='flex items-center gap-3 text-xl font-bold text-gray-900'>
            <Receipt className='h-5 w-5' />
            Bill Details - {bill.billNumber}
          </DialogTitle>
        </DialogHeader>
        <div className='rounded-lg border border-gray-200 bg-white p-6'>
          <div className='flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between'>
            <div className='flex items-center gap-4'>
              <img
                src='https://www.s3corp.com.vn/images/S3CORP.svg?w=128&q=75'
                alt='Company Logo'
                className='h-12 w-auto object-contain'
              />
              <div>
                <p className='text-lg font-bold text-gray-900'>S3Corp.</p>
                <p className='text-sm text-gray-600'>
                  {' '}
                  307/12 Nguyen Van Troi, Ward 1, Tan Binh District, HCMC, Viet Nam
                </p>
                <p className='text-sm text-gray-600'>Email: info@s3corp.com.vn | (+84) 28 3547 1411</p>
              </div>
            </div>

            <div className='text-right text-sm text-gray-700'>
              <p className='font-medium'>Printed on:</p>
              <p>{formatDate(new Date().toISOString())}</p>
            </div>
          </div>
        </div>

        <div className='grid grid-cols-1 gap-6 lg:grid-cols-4'>
          <div className='space-y-6 lg:col-span-3'>
            <div className='rounded-lg border border-gray-200 bg-white p-6'>
              <h3 className='mb-4 flex items-center gap-2 text-lg font-semibold text-gray-900'>
                <FileText className='h-5 w-5' />
                General Information
              </h3>
              <div className='space-y-0'>
                <InfoRow
                  icon={Hash}
                  label='Bill Number'
                  value={bill.billNumber}
                  valueClassName='font-mono font-medium'
                />

                <InfoRow
                  icon={Calendar}
                  label='Created At'
                  value={formatDate(bill.createAt)}
                />
                <InfoRow
                  icon={Clock}
                  label='Last Updated'
                  value={formatDate(getLastUpdated())}
                />
                <InfoRow
                  icon={Shield}
                  label='Status'
                  value={bill.status}
                  badge={true}
                  badgeClassName={getStatusColor(bill.status)}
                />
                <InfoRow
                  icon={DollarSign}
                  label='Amount'
                  value={`$${getAssetCost().toLocaleString()}`}
                  valueClassName='font-bold text-green-600'
                />
                {bill.description && (
                  <InfoRow
                    icon={FileText}
                    label='Description'
                    value={bill.description}
                    valueClassName='whitespace-pre-wrap break-words'
                  />
                )}
              </div>
            </div>

            <div className='rounded-lg border border-gray-200 bg-white p-6'>
              <h3 className='mb-4 flex items-center gap-2 text-lg font-semibold text-gray-900'>
                <User className='h-5 w-5' />
                Creator Information
              </h3>
              <div className='flex items-center gap-4'>
                <Avatar className='h-12 w-12'>
                  <AvatarImage
                    src={bill.creator?.avatar}
                    alt={bill.creator?.fullName}
                  />
                  <AvatarFallback className='bg-gray-100 text-gray-900'>{getCreatorInitials()}</AvatarFallback>
                </Avatar>
                <div className='flex-1'>
                  <InfoRow
                    icon={User}
                    label='Full Name'
                    value={bill.creator?.fullName || `User ${bill.createdBy}`}
                    valueClassName='font-medium'
                  />
                  <InfoRow
                    icon={Mail}
                    label='Email'
                    value={bill.creator?.email || 'Email not available'}
                  />
                </div>
              </div>
            </div>

            <div className='rounded-lg border border-gray-200 bg-white p-6'>
              <h3 className='mb-4 flex items-center gap-2 text-lg font-semibold text-gray-900'>
                <Package className='h-5 w-5' />
                Asset Information
              </h3>
              <div className='space-y-0'>
                <InfoRow
                  icon={FileText}
                  label='Asset Name'
                  value={getAssetName()}
                  valueClassName='font-medium'
                />
                {bill.assets?.serialNumber && (
                  <InfoRow
                    icon={Hash}
                    label='Serial Number'
                    value={bill.assets.serialNumber}
                    valueClassName='font-mono'
                  />
                )}
                <InfoRow
                  icon={Tag}
                  label='Category'
                  value={getCategoryName()}
                />
                <InfoRow
                  icon={DollarSign}
                  label='Cost'
                  value={`$${getAssetCost().toLocaleString()}`}
                />
                {bill.assets?.department && (
                  <InfoRow
                    icon={Building2}
                    label='Department'
                    value={bill.assets.department.departmentName}
                  />
                )}
                {bill.assets?.status && (
                  <InfoRow
                    icon={Shield}
                    label='Asset Status'
                    value={bill.assets.status}
                    badge={true}
                    badgeClassName='bg-green-100 text-green-800'
                  />
                )}
                {bill.assets?.purchaseDate && (
                  <InfoRow
                    icon={Calendar}
                    label='Purchase Date'
                    value={formatDateShort(bill.assets.purchaseDate)}
                  />
                )}
                {bill.assets?.warrantyExpiry && (
                  <InfoRow
                    icon={Clock}
                    label='Warranty Expiry'
                    value={`${formatDateShort(bill.assets.warrantyExpiry)} ${isWarrantyExpired(bill.assets.warrantyExpiry) ? '(Expired)' : ''}`}
                    valueClassName={isWarrantyExpired(bill.assets.warrantyExpiry) ? 'text-red-600' : ''}
                  />
                )}
              </div>
            </div>

            {(bill.fileAttachment || bill.imageUpload) && (
              <div className='rounded-lg border border-gray-200 bg-white p-6'>
                <h3 className='mb-4 flex items-center gap-2 text-lg font-semibold text-gray-900'>
                  <Paperclip className='h-5 w-5' />
                  File đính kèm
                </h3>
                <div className='space-y-0'>
                  {bill.fileAttachment && (
                    <InfoRow
                      icon={FileText}
                      label='Document'
                      value={bill.fileAttachment}
                      valueClassName='truncate max-w-xs'
                    />
                  )}
                  {bill.imageUpload && (
                    <InfoRow
                      icon={ImageIcon}
                      label='Image'
                      value={bill.imageUpload}
                      valueClassName='truncate max-w-xs'
                    />
                  )}
                </div>
              </div>
            )}
          </div>

          <div className='lg:col-span-1'>
            <div className='sticky top-4 rounded-lg border border-gray-200 bg-white p-6'>
              <h3 className='mb-4 flex items-center gap-2 text-lg font-semibold text-gray-900'>
                <QrCode className='h-5 w-5' />
                QR Code
              </h3>
              <div className='text-center'>
                <BillQR bill={bill} />
                <div className='mt-3 space-y-1'>
                  <p className='text-xs text-gray-500'>{bill.billNumber}</p>
                  <a
                    href={`${window.location.origin}/bills/${bill.id}`}
                    target='_blank'
                    rel='noopener noreferrer'
                    className='inline-flex items-center gap-1 text-xs text-blue-600 hover:text-blue-800'
                  >
                    <ExternalLink className='h-3 w-3' />
                    View Details
                  </a>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div className='rounded-lg border border-gray-200 bg-white p-6'>
          <h3 className='mb-4 flex items-center gap-2 text-lg font-semibold text-gray-900'>
            <AlertCircle className='h-5 w-5 text-yellow-500' />
            Notes / Remarks
          </h3>
          <p className='text-sm leading-relaxed text-gray-700'>
            Please verify all information carefully. If there are any issues with the billing details, contact the
            finance department within 3 business days.
          </p>
        </div>
        <div className='rounded-lg border border-gray-200 bg-white p-6'>
          <div className='grid grid-cols-2 gap-10 pt-6 text-sm text-gray-700'>
            <div className='text-center'>
              <p className='mb-12'>Created by</p>
              <p className='font-semibold underline'>{bill.creator?.fullName || '....................'}</p>
            </div>
            <div className='text-center'>
              <p className='mb-12'>Approved by</p>
              <p className='font-semibold underline'>....................</p>
            </div>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}
