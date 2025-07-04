import { useParams, Link, useNavigate } from 'react-router-dom'
import type { AssetsType } from '../view-all-assets/model'
import { getData, tryCatch } from '@/utils'
import { useEffect } from 'react'
import { deleteAsset, getAssetInformation } from '../api'
import { useState } from 'react'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  Button,
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
  LoadingSpinner,
} from '@/components/ui'

import { ArrowLeft, Pencil, Loader2, Trash2 } from 'lucide-react'
import { useAppSelector } from '@/hooks'

import {
  AssetBadge,
  AssetFile,
  AssetImage,
  AssetInformation,
  AssetMaintenanceSchedule,
  AssetQR,
  NoAsset,
} from './_components'
import { ViewAssetLog } from '../view-asset-log'
import { toast } from 'sonner'
import { CompareSimilarAssets } from './_components/compare-similar-assets'
import { ConfirmDeleteModal } from '@/components/ui/confirm-delete-modal'

const GetAssetDetail = () => {
  const { id } = useParams()
  const [isPending, setIsPending] = useState(false)
  const [isDeleting, setIsDeleting] = useState(false)
  const [asset, setAsset] = useState<AssetsType>()
  const [showDeleteModal, setShowDeleteModal] = useState(false)

  const navigate = useNavigate()
  const user = useAppSelector((state) => state.auth.user)
  const role = user.role.slug
  const canEditAsset = () => {
    if (role === 'admin') return true
    if (role === 'assetManager') {
      return user.department && asset?.department?.id === user.department.id
    }
    return false
  }
  const getAssetData = async () => {
    setIsPending(true)
    if (!id) return
    await getData(() => getAssetInformation(id), setAsset)
    setIsPending(false)
  }
  const deletingAsset = async () => {
    setIsDeleting(true)
    if (!id) return
    const response = await tryCatch(deleteAsset(id))
    if (response.error) {
      toast.error(response.error.message || 'Failed to delete asset')
      setIsDeleting(false)
      setShowDeleteModal(false)
      return
    }
    toast.success('Asset deleted successfully')
    navigate('/assets')
    setIsDeleting(false)
    setShowDeleteModal(false)
  }

  const handleDeleteClick = () => {
    setShowDeleteModal(true)
  }

  const handleDeleteConfirm = () => {
    deletingAsset()
  }

  const handleDeleteCancel = () => {
    setShowDeleteModal(false)
  }

  useEffect(() => {
    getAssetData()
  }, [id])

  if (isPending) {
    return (
      <div className='flex h-[70vh] items-center justify-center'>
        <Loader2 className='text-primary h-12 w-12 animate-spin' />
      </div>
    )
  }
  if (!asset) {
    return <NoAsset id={id || ''} />
  }

  return (
    <div className='container mx-auto px-4 py-4 sm:py-8'>
      <div className='mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between'>
        <div className='flex items-center'>
          <Link to='/assets'>
            <Button
              variant='ghost'
              className='mr-2 sm:mr-4'
              onClick={() => navigate(-1)}
            >
              <ArrowLeft className='h-5 w-5' />
            </Button>
          </Link>
          <div className='flex flex-col sm:flex-row sm:items-center sm:gap-2'>
            <h1 className='text-xl font-semibold sm:text-2xl lg:text-3xl'>{asset.assetName}</h1>
            <AssetBadge asset={asset} />
          </div>
        </div>

        {asset.status !== 'Disposed' && canEditAsset() && (
          <div className='flex flex-col gap-2 sm:flex-row sm:items-center sm:space-x-2'>
            <CompareSimilarAssets
              currentAsset={asset}
              className='flex-1'
            />
            <Link
              to={`/assets/update/${id}`}
              className='flex-1'
            >
              <Button
                variant='outline'
                className='w-full'
              >
                <Pencil className='mr-1 h-4 w-4' />
                <span className='hidden sm:inline'>Update Asset</span>
                <span className='sm:hidden'>Update Asset</span>
              </Button>
            </Link>

            <Button
              variant='destructive'
              onClick={handleDeleteClick}
              disabled={isDeleting}
              className='flex-1'
            >
              {isDeleting ? (
                <LoadingSpinner className='mr-2 h-4 w-4 animate-spin' />
              ) : (
                <>
                  <Trash2 className='mr-1 h-4 w-4' />
                  <span className='hidden sm:inline'>Delete Asset</span>
                  <span className='sm:hidden'>Delete Asset</span>
                </>
              )}
            </Button>
          </div>
        )}
      </div>

      <div className='grid grid-cols-1 gap-4 sm:gap-6 xl:grid-cols-3'>
        <div className='flex flex-col gap-4 sm:gap-6 xl:col-span-2'>
          <Card className='w-full'>
            <CardHeader>
              <CardTitle className='text-lg sm:text-xl'>Asset Information</CardTitle>
              <CardDescription>Details about the hardware asset</CardDescription>
            </CardHeader>
            <CardContent className='space-y-4 sm:space-y-6'>
              <AssetInformation asset={asset} />
            </CardContent>
          </Card>

          {role !== 'viewer' &&
            (role === 'admin' ||
              (role === 'assetManager' && (user as any).department?.id === asset?.department?.id)) && (
              <ViewAssetLog id={id || ''} />
            )}
          {role !== 'viewer' && (
            <Card className='w-full'>
              <div className='flex items-center justify-between'>
                <CardHeader className='flex-1'>
                  <CardTitle className='text-lg sm:text-xl'>Maintenance Schedule</CardTitle>
                  <CardDescription>Upcoming and past maintenance schedules</CardDescription>
                </CardHeader>
              </div>

              <CardContent>
                <AssetMaintenanceSchedule id={id || ''} />
              </CardContent>
            </Card>
          )}
        </div>

        <div className='flex xl:col-span-1'>
          <Tabs
            defaultValue='image'
            className='flex w-full flex-col'
          >
            <TabsList className='grid w-full grid-cols-3'>
              <TabsTrigger
                value='image'
                className='text-xs sm:text-sm'
              >
                Image
              </TabsTrigger>
              <TabsTrigger
                value='documents'
                className='text-xs sm:text-sm'
              >
                <span className='hidden sm:inline'>Documents</span>
                <span className='sm:hidden'>Docs</span>
              </TabsTrigger>
              <TabsTrigger
                value='qr'
                className='text-xs sm:text-sm'
              >
                QR
              </TabsTrigger>
            </TabsList>

            <TabsContent
              value='image'
              className='mt-2 flex flex-grow sm:mt-4'
            >
              <Card className='w-full'>
                <CardContent className='flex h-full min-h-[200px] items-center justify-center py-4 sm:min-h-[300px]'>
                  <AssetImage asset={asset} />
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent
              value='documents'
              className='mt-2 flex flex-grow sm:mt-4'
            >
              <Card className='w-full'>
                <CardContent className='flex h-full min-h-[200px] items-center justify-center py-4 sm:min-h-[300px]'>
                  <AssetFile asset={asset} />
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent
              value='qr'
              className='mt-2 flex flex-grow sm:mt-4'
            >
              <Card className='w-full'>
                <CardContent className='flex h-full min-h-[200px] items-center justify-center py-4 sm:min-h-[300px]'>
                  <AssetQR asset={asset} />
                </CardContent>
              </Card>
            </TabsContent>
          </Tabs>
        </div>
      </div>
      <ConfirmDeleteModal
        isOpen={showDeleteModal}
        onClose={handleDeleteCancel}
        onConfirm={handleDeleteConfirm}
        title='Confirm Delete'
        itemName={asset?.assetName}
        isLoading={isDeleting}
      />
    </div>
  )
}

export default GetAssetDetail
