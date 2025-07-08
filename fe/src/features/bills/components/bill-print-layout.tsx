import { Avatar, AvatarFallback, AvatarImage, Badge } from '@/components/ui'
import {
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
  Package,
  Clock,
  AlertCircle,
  CheckCircle,
  Phone,
  MapPin,
} from 'lucide-react'
import type { BillType } from '../model/bill-types'
import { BillQR } from './bill-qr'

interface BillPrintLayoutProps {
  bill: BillType
}

export const BillPrintLayout = ({ bill }: BillPrintLayoutProps) => {
  const getStatusColor = (status: string) => {
    const colors = {
      Unpaid: 'bg-orange-100 text-orange-800 border-orange-200',
      Paid: 'bg-green-100 text-green-800 border-green-200',
      Cancelled: 'bg-red-100 text-red-800 border-red-200',
      Draft: 'bg-gray-100 text-gray-800 border-gray-200',
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
  console.log('🚀 ~ isWarrantyExpired ~ isWarrantyExpired:', isWarrantyExpired)

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
    <div className='flex items-start gap-3 border-b border-gray-100 py-2 last:border-b-0 print:border-gray-300'>
      <Icon className='mt-0.5 h-4 w-4 flex-shrink-0 text-gray-500 print:text-gray-700' />
      <div className='flex-1 text-sm text-gray-700 print:text-gray-900'>
        <span className='font-medium text-gray-800 print:text-black'>{label}: </span>
        {badge ? (
          <Badge
            className={`${badgeClassName} ml-1 print:border print:bg-gray-100 print:text-gray-900`}
            variant='outline'
          >
            {value}
          </Badge>
        ) : (
          <span className={`text-gray-900 print:text-black ${valueClassName}`}>{value}</span>
        )}
      </div>
    </div>
  )

  return (
    <div className='bill-print-container min-h-screen bg-white p-8 text-black'>
      <style>{`
        @media print {
          .bill-print-container {
            padding: 15px;
            background: white !important;
            color: black !important;
            font-size: 12px;
          }
          .no-print {
            display: none !important;
          }
          .print\\:border-gray-300 {
            border-color: #d1d5db !important;
          }
          .print\\:text-gray-900, .print\\:text-black {
            color: #111827 !important;
          }
          .print\\:text-gray-700 {
            color: #374151 !important;
          }
          .print\\:bg-gray-100 {
            background-color: #f3f4f6 !important;
          }
          .print\\:border {
            border: 1px solid #d1d5db !important;
          }
          .print-signature {
            border-top: 2px solid #000 !important;
            padding-top: 8px !important;
          }
          .page-break {
            page-break-after: always;
          }
        }
        
        @page {
          margin: 1cm;
          size: A4;
        }
      `}</style>

      <div className='mb-6 rounded-lg border-2 border-gray-300 p-6 print:border-gray-400'>
        <div className='flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between'>
          <div className='flex items-start gap-4'>
            <img
              src='https://www.s3corp.com.vn/images/S3CORP.svg?w=128&q=75'
              alt='Company Logo'
              className='h-16 w-auto object-contain print:h-12'
            />
            <div className='flex-1'>
              <h1 className='text-2xl font-bold text-gray-900 print:text-xl'>S3Corp.</h1>
              <div className='mt-2 space-y-1 text-sm text-gray-600 print:text-gray-800'>
                <div className='flex items-center gap-1'>
                  <MapPin className='h-3 w-3' />
                  <span>307/12 Nguyen Van Troi, Ward 1, Tan Binh District, HCMC, Vietnam</span>
                </div>
                <div className='flex items-center gap-1'>
                  <Mail className='h-3 w-3' />
                  <span>info@s3corp.com.vn</span>
                </div>
                <div className='flex items-center gap-1'>
                  <Phone className='h-3 w-3' />
                  <span>(+84) 28 3547 1411</span>
                </div>
              </div>
            </div>
          </div>

          <div className='text-right'>
            <h2 className='text-3xl font-bold text-gray-900 print:text-2xl'>BILL INVOICE</h2>
            <p className='mt-1 text-sm text-gray-600 print:text-gray-800'>
              Generated on {formatDate(new Date().toISOString())}
            </p>
            <div className='mt-2 text-lg font-semibold text-blue-600 print:text-black'>{bill.billNumber}</div>
          </div>
        </div>
      </div>

      <div className='relative grid grid-cols-1 gap-6 lg:grid-cols-4'>
        <div className='no-print lg:no-print-override lg:absolute lg:bottom-0 lg:left-0 lg:w-32 print:w-24'>
          <div className='rounded-lg border border-gray-200 bg-white p-3 print:border-gray-300'>
            <div className='text-center'>
              <BillQR bill={bill} />
              <p className='mt-1 truncate text-xs text-gray-500 print:text-gray-700'>{bill.billNumber}</p>
            </div>
          </div>
        </div>

        <div className='space-y-6 lg:col-span-4 lg:pr-36'>
          <div className='rounded-lg border border-gray-200 p-6 print:border-gray-300'>
            <h3 className='mb-4 flex items-center gap-2 text-lg font-semibold text-gray-900 print:text-base'>
              <FileText className='h-5 w-5 print:h-4 print:w-4' />
              General Information
            </h3>
            <div className='grid grid-cols-1 gap-6 md:grid-cols-2'>
              <div className='space-y-0'>
                <InfoRow
                  icon={Hash}
                  label='Bill Number'
                  value={bill.billNumber}
                  valueClassName='font-mono font-bold'
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
              </div>
              <div className='space-y-0'>
                <InfoRow
                  icon={Shield}
                  label='Status'
                  value={bill.status || 'Unpaid'}
                  badge={true}
                  badgeClassName={getStatusColor(bill.status || 'Unpaid')}
                />
                <InfoRow
                  icon={DollarSign}
                  label='Amount'
                  value={`$${getAssetCost().toLocaleString()}`}
                  valueClassName='font-bold text-green-600 print:text-black text-lg'
                />
              </div>
            </div>
            {bill.description && (
              <div className='mt-4 border-t border-gray-100 pt-4 print:border-gray-300'>
                <InfoRow
                  icon={FileText}
                  label='Description'
                  value={bill.description}
                  valueClassName='whitespace-pre-wrap break-words italic'
                />
              </div>
            )}
          </div>

          <div className='grid grid-cols-1 gap-6 md:grid-cols-2'>
            <div className='rounded-lg border border-gray-200 p-6 print:border-gray-300'>
              <h3 className='mb-4 flex items-center gap-2 text-lg font-semibold text-gray-900 print:text-base'>
                <User className='h-5 w-5 print:h-4 print:w-4' />
                Creator Information
              </h3>
              <div className='mb-4 flex items-center gap-4'>
                <Avatar className='h-12 w-12 print:h-8 print:w-8'>
                  <AvatarImage
                    src={bill.creator?.avatar}
                    alt={bill.creator?.fullName}
                  />
                  <AvatarFallback className='bg-gray-100 text-gray-900 print:text-xs'>
                    {getCreatorInitials()}
                  </AvatarFallback>
                </Avatar>
                <div className='flex-1'>
                  <p className='font-semibold text-gray-900 print:text-sm'>
                    {bill.creator?.fullName || `User ${bill.createdBy}`}
                  </p>
                  <p className='text-sm text-gray-600 print:text-xs print:text-gray-800'>
                    {bill.creator?.email || 'Email not available'}
                  </p>
                </div>
              </div>
            </div>

            <div className='rounded-lg border border-gray-200 p-6 print:border-gray-300'>
              <h3 className='mb-4 flex items-center gap-2 text-lg font-semibold text-gray-900 print:text-base'>
                <Package className='h-5 w-5 print:h-4 print:w-4' />
                Asset Information
              </h3>
              <div className='space-y-0'>
                <InfoRow
                  icon={FileText}
                  label='Asset Name'
                  value={getAssetName()}
                  valueClassName='font-semibold'
                />
                {bill.assets?.serialNumber && (
                  <InfoRow
                    icon={Hash}
                    label='Serial Number'
                    value={bill.assets.serialNumber}
                    valueClassName='font-mono text-sm'
                  />
                )}
                <InfoRow
                  icon={Tag}
                  label='Category'
                  value={getCategoryName()}
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
                    badgeClassName='bg-blue-100 text-blue-800 border-blue-200'
                  />
                )}
              </div>
            </div>
          </div>

          {(bill.fileAttachment || bill.imageUpload) && (
            <div className='rounded-lg border border-gray-200 p-6 print:border-gray-300'>
              <h3 className='mb-4 flex items-center gap-2 text-lg font-semibold text-gray-900 print:text-base'>
                <Paperclip className='h-5 w-5 print:h-4 print:w-4' />
                Attachments
              </h3>
              <div className='space-y-0'>
                {bill.fileAttachment && (
                  <InfoRow
                    icon={FileText}
                    label='Document'
                    value={bill.fileAttachment}
                    valueClassName='break-all text-blue-600 print:text-black'
                  />
                )}
                {bill.imageUpload && (
                  <InfoRow
                    icon={ImageIcon}
                    label='Image'
                    value={bill.imageUpload}
                    valueClassName='break-all text-blue-600 print:text-black'
                  />
                )}
              </div>
            </div>
          )}

          <div className='rounded-lg border border-gray-200 p-6 print:border-gray-300'>
            <h3 className='mb-4 flex items-center gap-2 text-lg font-semibold text-gray-900 print:text-base'>
              <AlertCircle className='h-5 w-5 text-amber-500 print:h-4 print:w-4' />
              Important Notes & Terms
            </h3>
            <div className='space-y-3 text-sm leading-relaxed text-gray-700 print:text-xs print:text-gray-900'>
              <div className='flex items-start gap-2'>
                <CheckCircle className='mt-0.5 h-4 w-4 flex-shrink-0 text-green-500 print:h-3 print:w-3' />
                <p>
                  <strong>Verification:</strong> Please verify all information carefully before processing payment. Any
                  discrepancies should be reported within 3 business days.
                </p>
              </div>
              <div className='flex items-start gap-2'>
                <CheckCircle className='mt-0.5 h-4 w-4 flex-shrink-0 text-green-500 print:h-3 print:w-3' />
                <p>
                  <strong>Payment Terms:</strong> Payment is due within 30 days from bill generation date. Late payments
                  may incur additional fees.
                </p>
              </div>
              <div className='flex items-start gap-2'>
                <CheckCircle className='mt-0.5 h-4 w-4 flex-shrink-0 text-green-500 print:h-3 print:w-3' />
                <p>
                  <strong>Contact Information:</strong> For any billing inquiries, please contact our finance department
                  at finance@s3corp.com.vn or call (+84) 28 3547 1411.
                </p>
              </div>
              <div className='flex items-start gap-2'>
                <CheckCircle className='mt-0.5 h-4 w-4 flex-shrink-0 text-green-500 print:h-3 print:w-3' />
                <p>
                  <strong>Warranty:</strong>{' '}
                  {bill.assets?.warrantyExpiry
                    ? `Asset warranty expires on ${formatDateShort(bill.assets.warrantyExpiry)}`
                    : 'Please refer to purchase documentation for warranty information'}
                  .
                </p>
              </div>
              <div className='flex items-start gap-2'>
                <CheckCircle className='mt-0.5 h-4 w-4 flex-shrink-0 text-green-500 print:h-3 print:w-3' />
                <p>
                  <strong>Record Keeping:</strong> This document serves as an official record. Please retain for your
                  financial records and future reference.
                </p>
              </div>
            </div>
          </div>

          <div className='rounded-lg border border-gray-200 p-6 print:border-gray-300'>
            <h3 className='mb-6 flex items-center gap-2 text-lg font-semibold text-gray-900 print:text-base'>
              <FileText className='h-5 w-5 print:h-4 print:w-4' />
              Authorization & Approval
            </h3>
            <div className='grid grid-cols-1 gap-8 text-center md:grid-cols-3'>
              <div className='space-y-4'>
                <p className='text-sm font-medium text-gray-700 print:text-xs print:text-gray-900'>Prepared by</p>
                <div className='flex h-16 items-end print:h-12'>
                  <div className='w-full'>
                    <div className='print-signature border-t-2 border-gray-800 pt-2'>
                      <p className='text-sm font-semibold print:text-xs'>{bill.creator?.fullName || 'N/A'}</p>
                      <p className='mt-1 text-xs text-gray-600 print:text-gray-800'>Creator</p>
                    </div>
                  </div>
                </div>
                <p className='text-xs text-gray-500 print:text-gray-700'>Date: {formatDateShort(bill.createAt)}</p>
              </div>

              <div className='space-y-4'>
                <p className='text-sm font-medium text-gray-700 print:text-xs print:text-gray-900'>Reviewed by</p>
                <div className='flex h-16 items-end print:h-12'>
                  <div className='w-full'>
                    <div className='print-signature border-t-2 border-gray-800 pt-2'>
                      <p className='text-sm font-semibold print:text-xs'>_________________</p>
                      <p className='mt-1 text-xs text-gray-600 print:text-gray-800'>Finance Manager</p>
                    </div>
                  </div>
                </div>
                <p className='text-xs text-gray-500 print:text-gray-700'>Date: _______________</p>
              </div>

              <div className='space-y-4'>
                <p className='text-sm font-medium text-gray-700 print:text-xs print:text-gray-900'>Approved by</p>
                <div className='flex h-16 items-end print:h-12'>
                  <div className='w-full'>
                    <div className='print-signature border-t-2 border-gray-800 pt-2'>
                      <p className='text-sm font-semibold print:text-xs'>_________________</p>
                      <p className='mt-1 text-xs text-gray-600 print:text-gray-800'>Director</p>
                    </div>
                  </div>
                </div>
                <p className='text-xs text-gray-500 print:text-gray-700'>Date: _______________</p>
              </div>
            </div>
          </div>

          <div className='border-t border-gray-200 py-4 text-center print:border-gray-300'>
            <p className='text-xs text-gray-500 print:text-gray-700'>
              This is a computer-generated document. No signature is required.
            </p>
            <p className='mt-1 text-xs text-gray-500 print:text-gray-700'>© 2024 S3Corp. All rights reserved.</p>
          </div>
        </div>
      </div>
    </div>
  )
}
