import { useEffect, useState, useCallback, useMemo } from 'react'
import { useSearchParams, useNavigate } from 'react-router-dom'
import { useAppSelector, useDebounce } from '@/hooks'
import { MonthlyReport } from '@/features/dashboard/components/monthly-report'
import { getDataAssetsFilter, getAllAssets } from '@/features/assets/api'
import type { AssetsType } from '@/features/assets/view-all-assets/model'
import type { DateFilter } from '@/features/dashboard/model/statistics-types'
import type { FilterType } from '@/features/assets/view-all-assets/model'
import { tryCatch } from '@/utils'
import { AreaChart, FileText, Loader2 } from 'lucide-react'
import { Button } from '@/components/ui'

const StatisticReportPage = () => {
  const [searchParam, setSearchParam] = useSearchParams()
  const navigate = useNavigate()

  const [assets, setAssets] = useState<AssetsType[]>([])
  const [isInitialLoading, setIsInitialLoading] = useState(true)

  const [assetFilter, setAssetFilter] = useState<FilterType>({
    assetName: '',
    categoryId: searchParam.get('categoryId') || null,
    departmentId: searchParam.get('departmentId') || null,
    status: searchParam.get('status') || null,
  })

  const [dateFilter, setDateFilter] = useState<DateFilter>({
    dateField: 'purchase',
    month: searchParam.get('month') ? parseInt(searchParam.get('month')!) : undefined,
    year: searchParam.get('year') ? parseInt(searchParam.get('year')!) : undefined,
    startDate: searchParam.get('startDate') ? new Date(searchParam.get('startDate')!) : undefined,
    endDate: searchParam.get('endDate') ? new Date(searchParam.get('endDate')!) : undefined,
  })

  const assetFilterData = useDebounce(assetFilter, 300)

  const user = useAppSelector((state) => state.auth.user)
  const role = user.role.slug
  const isAdmin = role === 'admin'

  const filterParams = useMemo(
    () => ({
      assetName: '',
      categoryId: assetFilterData.categoryId,
      departmentId: assetFilterData.departmentId,
      status: assetFilterData.status,
      month: dateFilter.month,
      year: dateFilter.year,
      startDate: dateFilter.startDate?.toISOString(),
      endDate: dateFilter.endDate?.toISOString(),
      dateField: dateFilter.dateField,
    }),
    [assetFilterData, dateFilter]
  )

  const hasActiveFilters = useMemo(() => {
    return Object.values(filterParams).some((value) => value !== null && value !== '')
  }, [filterParams])

  const fetchAssets = useCallback(
    async (isInitial = false) => {
      if (isInitial) setIsInitialLoading(true)

      try {
        const response = hasActiveFilters
          ? await tryCatch(getDataAssetsFilter(filterParams))
          : await tryCatch(getAllAssets())

        if (!response.error) {
          setAssets(response.data.data || [])
        }
      } catch (error) {
        console.error('Fetch error:', error)
      } finally {
        if (isInitial) setIsInitialLoading(false)
      }
    },
    [hasActiveFilters, filterParams]
  )

  useEffect(() => {
    fetchAssets(true)
  }, [])

  useEffect(() => {
    if (!isInitialLoading) {
      fetchAssets(false)
    }
  }, [filterParams, isInitialLoading])

  useEffect(() => {
    const params = new URLSearchParams()

    if (dateFilter.month) params.set('month', dateFilter.month.toString())
    if (dateFilter.year) params.set('year', dateFilter.year.toString())
    if (dateFilter.startDate) params.set('startDate', dateFilter.startDate.toISOString())
    if (dateFilter.endDate) params.set('endDate', dateFilter.endDate.toISOString())

    if (assetFilterData.categoryId) params.set('categoryId', assetFilterData.categoryId)
    if (assetFilterData.departmentId) params.set('departmentId', assetFilterData.departmentId)
    if (assetFilterData.status) params.set('status', assetFilterData.status)

    if (params.toString() !== searchParam.toString()) {
      setSearchParam(params, { replace: true })
    }
  }, [dateFilter, assetFilterData, searchParam, setSearchParam])

  const handleAssetFilterChange = useCallback((newFilter: FilterType) => {
    setAssetFilter(newFilter)
  }, [])

  const handleDateFilterChange = useCallback((newFilter: DateFilter) => {
    setDateFilter(newFilter)
  }, [])

  const handleNavigateToDashboard = useCallback(() => {
    navigate('/dashboard')
  }, [navigate])

  const views = useMemo(
    () => [
      {
        key: 'dashboard',
        label: 'Dashboard',
        icon: AreaChart,
        onClick: handleNavigateToDashboard,
      },
      ...(isAdmin
        ? [
            {
              key: 'reports',
              label: 'Statistical Reports',
              icon: FileText,
              onClick: () => {},
            },
          ]
        : []),
    ],
    [isAdmin, handleNavigateToDashboard]
  )

  if (isInitialLoading) {
    return (
      <div className='flex min-h-screen items-center justify-center'>
        <div className='text-center'>
          <Loader2 className='mx-auto h-8 w-8 animate-spin dark:text-gray-300' />
          <p className='mt-2 text-gray-600 dark:text-gray-400'>Loading assets...</p>
        </div>
      </div>
    )
  }

  return (
    <div className='space-y-4 p-4 sm:space-y-6 sm:p-6'>
      <div className='flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between'>
        <h1 className='flex items-center gap-2 text-2xl font-bold sm:text-3xl'>
          <FileText className='h-6 w-6 sm:h-8 sm:w-8' />
          Statistical Reports
        </h1>

        {views.length > 1 && (
          <div className='flex gap-2'>
            {views.map((view) => (
              <Button
                key={view.key}
                variant={view.key === 'reports' ? 'default' : 'outline'}
                onClick={view.onClick}
                type='button'
                className='flex items-center gap-2'
              >
                <view.icon className='h-4 w-4' />
                {view.label}
              </Button>
            ))}
          </div>
        )}
      </div>

      <MonthlyReport
        assets={assets}
        initialDateFilter={dateFilter}
        onDateFilterChange={handleDateFilterChange}
        assetFilter={assetFilter}
        onAssetFilterChange={handleAssetFilterChange}
      />
    </div>
  )
}

export default StatisticReportPage
