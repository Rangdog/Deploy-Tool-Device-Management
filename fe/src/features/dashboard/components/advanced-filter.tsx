import { Card, CardContent, CardHeader, CardTitle, Button } from '@/components/ui'
import { Filter, X } from 'lucide-react'
import type { FilterType } from '../../assets/view-all-assets/model'
import type { AssetsType } from '../../assets/view-all-assets/model'
import type { DateFilter } from '../model/statistics-types'
import type { SetStateAction } from 'react'
import { DateFilter as DateFilterComponent } from './date-filter'
import {
  CategoriesFilter,
  DepartmentsFilter,
  StatusFilter,
} from '../../assets/view-all-assets/_components/filter/_components'
import { useState } from 'react'

interface AdvancedFilterProps {
  dateFilter: DateFilter
  onDateFilterChange: (filter: DateFilter) => void
  assetFilter: FilterType
  onAssetFilterChange: (filter: FilterType) => void
  assets: AssetsType[]
  className?: string
}

export const AdvancedFilter = ({
  dateFilter,
  onDateFilterChange,
  assetFilter,
  onAssetFilterChange,
  assets,
  className = '',
}: AdvancedFilterProps) => {
  const handleAssetFilterChange: React.Dispatch<SetStateAction<FilterType>> = (value) => {
    const newFilter = typeof value === 'function' ? value(assetFilter) : value
    onAssetFilterChange(newFilter)
  }
  const [originalAssets] = useState(assets)
  const handleClearAllFilters = () => {
    onDateFilterChange({
      dateField: 'purchase',
      month: undefined,
      year: undefined,
      startDate: undefined,
      endDate: undefined,
      singleDate: undefined,
    })

    onAssetFilterChange({
      assetName: '',
      categoryId: null,
      departmentId: null,
      status: null,
    })
  }

  const hasActiveFilters =
    dateFilter.month ||
    dateFilter.year ||
    dateFilter.startDate ||
    assetFilter.categoryId ||
    assetFilter.departmentId ||
    assetFilter.status

  return (
    <Card className={className}>
      <CardHeader>
        <div className='flex items-center justify-between'>
          <CardTitle className='flex items-center gap-2'>
            <Filter className='h-5 w-5' />
            Filters
          </CardTitle>
          {hasActiveFilters && (
            <Button
              variant='outline'
              size='sm'
              onClick={handleClearAllFilters}
              className='gap-2'
            >
              <X className='h-4 w-4' />
              Clear All
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent className='space-y-4'>
        <div className='flex flex-wrap items-end gap-4'>
          <DateFilterComponent
            dateFilter={dateFilter}
            onDateFilterChange={onDateFilterChange}
            originalAssets={originalAssets}
            assets={assets}
          />

          <div className='space-y-2'>
            <label className='text-sm font-medium text-gray-600'>Category</label>
            <div className='h-10'>
              <CategoriesFilter
                filteredAssets={assetFilter}
                setFilteredAssets={handleAssetFilterChange}
              />
            </div>
          </div>

          <div className='space-y-2'>
            <label className='text-sm font-medium text-gray-600'>Department</label>
            <div className='h-10'>
              <DepartmentsFilter
                filteredAssets={assetFilter}
                setFilteredAssets={handleAssetFilterChange}
              />
            </div>
          </div>

          <div className='space-y-2'>
            <label className='text-sm font-medium text-gray-600'>Status</label>
            <div className='h-10'>
              <StatusFilter
                filteredAssets={assetFilter}
                setFilteredAssets={handleAssetFilterChange}
              />
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
