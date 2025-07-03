import { useState, useEffect, useRef } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogFooter,
  Button,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Input,
  Label,
  Badge,
} from '@/components/ui'
import { Plus, Receipt, File, Image, X, Shield } from 'lucide-react'
import { createBill } from '../api/create-bill'
import { getAllAssets } from '@/features/assets/api/get-all-assets'
import { createBillSchema, type CreateBillFormType } from '../model/create-bill-schema'
import type { BillType, CreateBillRequest } from '../model/bill-types'
import { tryCatch } from '@/utils'
import { toast } from 'sonner'
import Cookies from 'js-cookie'
import { useAppSelector } from '@/hooks'

interface CreateBillModalProps {
  onBillCreated: (bill: BillType) => void
}

interface AssetWithCategory {
  id: string
  assetName: string
  cost: number
  status: string
  category: {
    id: number
    categoryName: string
  }
}

export const CreateBillModal = ({ onBillCreated }: CreateBillModalProps) => {
  const [open, setOpen] = useState(false)
  const [assets, setAssets] = useState<AssetWithCategory[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [selectedFiles, setSelectedFiles] = useState<{
    fileAttachment?: File
    imageUpload?: File
  }>({})

  const fileInputRef = useRef<HTMLInputElement>(null)
  const imageInputRef = useRef<HTMLInputElement>(null)

  const {
    register,
    handleSubmit,
    formState: { errors },
    setValue,
    watch,
    reset,
  } = useForm<CreateBillFormType>({
    resolver: zodResolver(createBillSchema),
    defaultValues: {
      assetId: '',
      assetName: '',
      categoryName: '',
      cost: 0,
      description: '',
      status: 'Unpaid',
      fileAttachment: '',
      imageUpload: '',
    },
  })

  const currentUser = useAppSelector((state) => state.auth.user)

  const getCurrentUserId = () => {
    const id = Cookies.get('id') || currentUser.id
    return parseInt(id) || parseInt(currentUser.id) || 1
  }

  useEffect(() => {
    if (open) {
      const fetchData = async () => {
        setIsLoading(true)
        try {
          const response = await tryCatch(getAllAssets())
          console.log('Assets response:', response)

          if (!response.error && response.data?.data) {
            const mappedAssets = response.data.data.map((asset: any) => ({
              id: asset.id.toString(),
              assetName: asset.assetName,
              cost: asset.cost || 0,
              status: asset.status || 'Active',
              category: asset.category || { id: 0, categoryName: 'No Category' },
            }))
            setAssets(mappedAssets)
          } else {
            console.error('Failed to fetch assets:', response.error)
            toast.error('Failed to fetch assets')
          }
        } catch (error) {
          console.error('Assets fetch error:', error)
          toast.error('Failed to fetch assets')
        }
        setIsLoading(false)
      }
      fetchData()
    }
  }, [open])

  const statusOptions = [
    { value: 'Unpaid', label: 'Unpaid' },
    { value: 'Paid', label: 'Paid' },
  ]

  const selectedAsset = watch('assetId')
  useEffect(() => {
    if (selectedAsset) {
      const asset = assets.find((a) => a.id === selectedAsset)
      if (asset) {
        setValue('assetName', asset.assetName)
        setValue('categoryName', asset.category.categoryName)
        setValue('cost', asset.cost)
      }
    } else {
      setValue('assetName', '')
      setValue('categoryName', '')
      setValue('cost', 0)
    }
  }, [selectedAsset, assets, setValue])

  const selectedAssetData = selectedAsset ? assets.find((a) => a.id === selectedAsset) : null

  const handleFileUpload = (type: 'fileAttachment' | 'imageUpload', file: File | null) => {
    if (file) {
      setSelectedFiles((prev) => ({ ...prev, [type]: file }))
      setValue(type, file)
    } else {
      setSelectedFiles((prev) => {
        const newFiles = { ...prev }
        delete newFiles[type]
        return newFiles
      })
      setValue(type, '')
    }
  }

  const removeFile = (type: 'fileAttachment' | 'imageUpload') => {
    handleFileUpload(type, null)
    if (type === 'fileAttachment' && fileInputRef.current) {
      fileInputRef.current.value = ''
    }
    if (type === 'imageUpload' && imageInputRef.current) {
      imageInputRef.current.value = ''
    }
  }

  const getAssetStatusColor = (status: string) => {
    const colors = {
      Active: 'bg-green-100 text-green-800',
      Inactive: 'bg-gray-100 text-gray-800',
      Maintenance: 'bg-yellow-100 text-yellow-800',
      Retired: 'bg-red-100 text-red-800',
      New: 'bg-blue-100 text-blue-800',
      'In Use': 'bg-green-100 text-green-800',
    } as const
    return colors[status as keyof typeof colors] || colors.Active
  }

  const onSubmit = async (data: CreateBillFormType) => {
    setIsSubmitting(true)

    try {
      const selectedAssetData = assets.find((a) => a.id === data.assetId)

      if (!selectedAssetData) {
        toast.error('Please select a valid asset')
        setIsSubmitting(false)
        return
      }

      const currentUserId = getCurrentUserId()

      const formData = new FormData()
      formData.append('assetId', data.assetId)
      formData.append('amount', data.cost.toString())
      formData.append('description', data.description.trim())
      formData.append('status', data.status)
      formData.append('categoryId', selectedAssetData.category.id.toString())
      formData.append('createdBy', currentUserId.toString())

      if (selectedFiles.fileAttachment) {
        formData.append('fileAttachment', selectedFiles.fileAttachment)
      }
      if (selectedFiles.imageUpload) {
        formData.append('imageUpload', selectedFiles.imageUpload)
      }

      const response = await tryCatch(createBill(formData as unknown as CreateBillRequest))

      if (response.error) {
        console.error('Create bill error:', response.error)
        toast.error(response.error?.message || 'Failed to create bill')
      } else {
        const newBill: BillType = {
          id: response.data?.data?.id || Date.now(),
          billNumber: response.data?.data?.billNumber || `BILL-${Date.now()}`,
          assetId: parseInt(data.assetId),
          amount: data.cost,
          description: data.description,
          status: data.status,
          categoryId: selectedAssetData.category.id,
          companyId: response.data?.data?.companyId || 1,
          createdBy: response.data?.data?.createdBy || 1,
          createAt: response.data?.data?.createdAt || new Date().toISOString(),
          updateAt: response.data?.data?.updatedAt || new Date().toISOString(),
          fileAttachment: response.data?.data?.fileAttachment,
          imageUpload: response.data?.data?.imageUpload,

          assets: {
            id: parseInt(selectedAssetData.id),
            assetName: data.assetName,
            cost: data.cost,
            status: selectedAssetData.status,
            category: {
              id: selectedAssetData.category.id,
              categoryName: data.categoryName,
            },
          },
        }

        onBillCreated(newBill)
        setOpen(false)
        reset()
        setSelectedFiles({})
        toast.success('Bill created successfully')
      }
    } catch (error) {
      console.error('Unexpected error:', error)
      toast.error('An unexpected error occurred')
    }

    setIsSubmitting(false)
  }

  const handleCancel = () => {
    reset()
    setSelectedFiles({})
    setOpen(false)
  }

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen && !isSubmitting) {
      reset()
      setSelectedFiles({})
      setOpen(false)
    }
  }

  return (
    <Dialog
      open={open}
      onOpenChange={handleOpenChange}
    >
      <DialogTrigger asChild>
        <Button
          className='flex items-center gap-2'
          onClick={() => setOpen(true)}
        >
          <Plus className='h-4 w-4' />
          Create Bill
        </Button>
      </DialogTrigger>

      <DialogContent className='max-h-[90vh] !max-w-2xl overflow-y-auto'>
        <DialogHeader>
          <DialogTitle className='flex items-center text-xl'>
            <Receipt className='mr-2 h-5 w-5' />
            Create New Bill
          </DialogTitle>
        </DialogHeader>

        <form
          onSubmit={handleSubmit(onSubmit)}
          className='space-y-6'
        >
          <div className='grid grid-cols-1 gap-6 lg:grid-cols-2'>
            <div className='space-y-6'>
              <div className='space-y-2'>
                <Label htmlFor='assetId'>Select Asset *</Label>
                <Select
                  value={watch('assetId')}
                  onValueChange={(value) => setValue('assetId', value)}
                  disabled={isLoading}
                >
                  <SelectTrigger className='h-11 w-full'>
                    <SelectValue placeholder='Choose an asset to create bill' />
                  </SelectTrigger>
                  <SelectContent>
                    {assets.map((asset) => (
                      <SelectItem
                        key={asset.id}
                        value={asset.id}
                      >
                        {asset.assetName}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {errors.assetId && <p className='text-sm text-red-500'>{errors.assetId.message}</p>}
              </div>

              {selectedAssetData && (
                <div className='space-y-3'>
                  <Label>Asset Information</Label>
                  <div className='space-y-3 rounded-lg border border-gray-200 bg-gray-50 p-4'>
                    <div className='flex items-center justify-between border-b border-gray-200 py-2 last:border-b-0'>
                      <span className='text-sm font-medium text-gray-700'>Asset Name:</span>
                      <span className='text-sm font-medium text-gray-900'>{watch('assetName')}</span>
                    </div>

                    <div className='flex items-center justify-between border-b border-gray-200 py-2 last:border-b-0'>
                      <span className='text-sm font-medium text-gray-700'>Category:</span>
                      <span className='text-sm text-gray-900'>{watch('categoryName')}</span>
                    </div>

                    <div className='flex items-center justify-between border-b border-gray-200 py-2 last:border-b-0'>
                      <span className='text-sm font-medium text-gray-700'>Cost:</span>
                      <span className='text-sm font-semibold text-green-600'>${watch('cost').toLocaleString()}</span>
                    </div>

                    <div className='flex items-center justify-between py-2'>
                      <span className='text-sm font-medium text-gray-700'>Status:</span>
                      <Badge
                        className={getAssetStatusColor(selectedAssetData.status)}
                        variant='outline'
                      >
                        <Shield className='mr-1 h-3 w-3' />
                        {selectedAssetData.status}
                      </Badge>
                    </div>
                  </div>
                </div>
              )}

              <div className='space-y-2'>
                <Label htmlFor='status'>Bill Status</Label>
                <Select
                  value={watch('status')}
                  onValueChange={(value) => setValue('status', value as any)}
                >
                  <SelectTrigger className='h-11 w-full'>
                    <SelectValue placeholder='Select status' />
                  </SelectTrigger>
                  <SelectContent>
                    {statusOptions.map((option) => (
                      <SelectItem
                        key={option.value}
                        value={option.value}
                      >
                        {option.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {errors.status && <p className='text-sm text-red-500'>{errors.status.message}</p>}
              </div>
            </div>

            <div className='space-y-6'>
              <div className='space-y-2'>
                <Label htmlFor='description'>Description *</Label>
                <Input
                  id='description'
                  placeholder='Enter bill description...'
                  className='h-9'
                  {...register('description')}
                />
                {errors.description && <p className='text-sm text-red-500'>{errors.description.message}</p>}
              </div>

              <div className='space-y-4'>
                <Label>File Attachments (Optional)</Label>

                <div className='space-y-2'>
                  <div className='flex items-center gap-2'>
                    <File className='h-4 w-4 text-gray-500' />
                    <span className='text-sm font-medium text-gray-700'>Document</span>
                  </div>

                  <div className='flex items-center gap-2'>
                    <Input
                      ref={fileInputRef}
                      type='file'
                      accept='.pdf,.doc,.docx,.txt'
                      onChange={(e) => {
                        const file = e.target.files?.[0]
                        handleFileUpload('fileAttachment', file || null)
                      }}
                      className='flex-1'
                    />
                    {selectedFiles.fileAttachment && (
                      <Button
                        type='button'
                        variant='outline'
                        size='sm'
                        onClick={() => removeFile('fileAttachment')}
                      >
                        <X className='h-4 w-4' />
                      </Button>
                    )}
                  </div>

                  {selectedFiles.fileAttachment && (
                    <p className='text-xs text-gray-600'>Selected: {selectedFiles.fileAttachment.name}</p>
                  )}
                </div>

                <div className='space-y-2'>
                  <div className='flex items-center gap-2'>
                    <Image className='h-4 w-4 text-gray-500' />
                    <span className='text-sm font-medium text-gray-700'>Image</span>
                  </div>

                  <div className='flex items-center gap-2'>
                    <Input
                      ref={imageInputRef}
                      type='file'
                      accept='image/*'
                      onChange={(e) => {
                        const file = e.target.files?.[0]
                        handleFileUpload('imageUpload', file || null)
                      }}
                      className='flex-1'
                    />
                    {selectedFiles.imageUpload && (
                      <Button
                        type='button'
                        variant='outline'
                        size='sm'
                        onClick={() => removeFile('imageUpload')}
                      >
                        <X className='h-4 w-4' />
                      </Button>
                    )}
                  </div>

                  {selectedFiles.imageUpload && (
                    <p className='text-xs text-gray-600'>Selected: {selectedFiles.imageUpload.name}</p>
                  )}
                </div>
              </div>
            </div>
          </div>

          {/* Hidden fields for validation */}
          <input
            type='hidden'
            {...register('assetName')}
          />
          <input
            type='hidden'
            {...register('categoryName')}
          />
          <input
            type='hidden'
            {...register('cost', { valueAsNumber: true })}
          />

          <DialogFooter className='flex gap-2'>
            <Button
              type='button'
              variant='outline'
              onClick={handleCancel}
              disabled={isSubmitting}
            >
              Cancel
            </Button>
            <Button
              type='submit'
              disabled={isSubmitting || isLoading || !selectedAssetData}
              className='flex items-center gap-2'
            >
              <Receipt className='h-4 w-4' />
              {isSubmitting ? 'Creating...' : 'Create Bill'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
