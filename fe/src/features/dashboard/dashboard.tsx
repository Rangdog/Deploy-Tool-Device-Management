import { useEffect, useState } from 'react'
import { useAppSelector, useDebounce } from '@/hooks'
import type { FilterType, AssetsType } from '../assets/view-all-assets'
import { DashboardFilter } from './components/dashboard-filter'
import { DashboardStats } from './components/dashboard-stats'
import { getAllAssets, getDataAssetsFilter } from '../assets/api'
import { tryCatch } from '@/utils'
import type { DashboardData } from './api/type'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { AreaChart, FileText } from 'lucide-react'
import { Button } from '@/components/ui'

export const Dashboard = () => {
  const [searchParam, setSearchParam] = useSearchParams()
  const navigate = useNavigate()
  const [filteredAssets, setFilteredAssets] = useState<FilterType>({
    assetName: searchParam.get('assetName') || '',
    categoryId: searchParam.get('categoryId') || null,
    departmentId: searchParam.get('departmentId') || null,
    status: searchParam.get('status') || null,
  })

  const [assets, setAssets] = useState<AssetsType[]>([])
  const [stats, setStats] = useState<DashboardData>({
    total_assets: 0,
    assigned: 0,
    under_maintenance: 0,
    retired: 0,
  })
  const [isPending, setIsPending] = useState(false)

  const filterData = useDebounce(filteredAssets, 1000)
  const user = useAppSelector((state) => state.auth.user)
  const role = user.role.slug
  const isAdmin = role === 'admin'

  const calculateStats = (assets: AssetsType[]) => {
    const newStats: DashboardData = {
      total_assets: assets.length,
      assigned: assets.filter((asset) => asset.status === 'In Use').length,
      under_maintenance: assets.filter((asset) => asset.status === 'Under Maintenance').length,
      retired: assets.filter((asset) => asset.status === 'Retired').length,
    }
    setStats(newStats)
  }

  const fetchAssets = async () => {
    setIsPending(true)
    try {
      if (Object.values(filterData).some((value) => value !== null && value !== '')) {
        const response = await tryCatch(getDataAssetsFilter(filterData))
        if (!response.error) {
          setAssets(response.data.data)
          calculateStats(response.data.data)
        }
      } else {
        const response = await tryCatch(getAllAssets())
        if (!response.error) {
          setAssets(response.data.data)
          calculateStats(response.data.data)
        }
      }
    } finally {
      setIsPending(false)
    }
  }

  useEffect(() => {
    if (filteredAssets.assetName) {
      searchParam.set('assetName', filteredAssets.assetName)
    } else {
      searchParam.delete('assetName')
    }
    if (filteredAssets.categoryId) {
      searchParam.set('categoryId', filteredAssets.categoryId)
    } else {
      searchParam.delete('categoryId')
    }
    if (filteredAssets.departmentId) {
      searchParam.set('departmentId', filteredAssets.departmentId)
    } else {
      searchParam.delete('departmentId')
    }
    if (filteredAssets.status) {
      searchParam.set('status', filteredAssets.status)
    } else {
      searchParam.delete('status')
    }
    setSearchParam(searchParam)
  }, [filterData])

  useEffect(() => {
    fetchAssets()
  }, [filterData])
  const handleNavigateToReports = () => {
    navigate('/statistics-reports')
  }
  const views = [
    {
      key: 'dashboard',
      label: 'Dashboard',
      icon: AreaChart,
      onClick: () => {},
    },
    ...(isAdmin
      ? [
          {
            key: 'reports',
            label: 'Statistical Reports',
            icon: FileText,
            onClick: handleNavigateToReports,
          },
        ]
      : []),
  ]

  return (
    <div className='space-y-4 p-4 sm:space-y-6 sm:p-6'>
      <div className='flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between'>
        <h1 className='flex items-center gap-2 text-2xl font-bold sm:text-3xl'>
          <AreaChart className='h-6 w-6 sm:h-8 sm:w-8' />
          Dashboard
        </h1>

        {views.length > 1 && (
          <div className='flex gap-2'>
            {views.map((view) => (
              <Button
                key={view.key}
                variant={view.key === 'dashboard' ? 'default' : 'outline'}
                onClick={view.onClick}
                className='flex items-center gap-2'
              >
                <view.icon className='h-4 w-4' />
                {view.label}
              </Button>
            ))}
          </div>
        )}
      </div>

      <DashboardFilter
        filteredAssets={filteredAssets}
        setFilteredAssets={setFilteredAssets}
        assets={assets}
      />

      <DashboardStats
        stats={stats}
        assets={assets}
        isPending={isPending}
      />
    </div>
  )
}
