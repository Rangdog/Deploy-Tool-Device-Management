import { useState, useEffect } from 'react'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  Button,
  Input,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui'
import { Filter, FilterX } from 'lucide-react'
import { getAllCategories } from '@/features/assets/api/get-all-categories'
import type { BillFilterType } from '../model/bill-types'
import { tryCatch } from '@/utils'
import { toast } from 'sonner'

interface BillsFilterProps {
  filters: BillFilterType
  setFilters: (filters: BillFilterType) => void
  onReset: () => void
}

export const BillsFilter = ({ filters, setFilters, onReset }: BillsFilterProps) => {
  const [categories, setCategories] = useState<Array<{ id: string; categoryName: string }>>([])
  const [isLoading, setIsLoading] = useState(false)

  useEffect(() => {
    const fetchData = async () => {
      setIsLoading(true)
      try {
        const response = await tryCatch(getAllCategories())
        if (!response.error && response.data?.data) {
          const mappedCategories = response.data.data.map((category: any) => ({
            id: category.id.toString(),
            categoryName: category.categoryName,
          }))
          setCategories(mappedCategories)
        } else {
          console.error('Failed to fetch categories:', response.error)
          toast.error('Failed to fetch categories')
        }
      } catch (error) {
        console.error('Categories fetch error:', error)
        toast.error('Failed to fetch categories')
      }
      setIsLoading(false)
    }
    fetchData()
  }, [])

  console.log('🚀 ~ mappedCategories ~ category:', categories)
  const statusOptions = [
    { id: 'Unpaid', statusName: 'Unpaid' },
    { id: 'Paid', statusName: 'Paid' },
  ]

  const hasActiveFilters = Object.values(filters).some((value) => value !== null && value !== '')

  console.log('🚀 ~ BillsFilter ~ statusOptions:', statusOptions)
  return (
    <Card>
      <CardHeader className='gap-0 pb-0'>
        <CardTitle className='flex items-center text-lg'>
          <Filter className='mr-2 h-4 w-4' />
          Filter Bills
        </CardTitle>
      </CardHeader>
      <CardContent className='pt-0'>
        <div className='flex flex-col gap-4 lg:flex-row lg:items-end lg:gap-3'>
          <div className='space-y-2'>
            <label className='text-sm font-medium'>Bill Number</label>
            <Input
              value={filters.billNumber || ''}
              onChange={(e) => setFilters({ ...filters, billNumber: e.target.value || '' })}
              placeholder='Search by bill number...'
              className='mt-2 h-9'
            />
          </div>

          <div className='flex flex-col gap-4 sm:flex-row sm:gap-2 lg:gap-3'>
            <div className='w-full space-y-2 sm:flex-1'>
              <label className='text-sm font-medium'>Category</label>
              <Select
                value={filters.categoryId || ''}
                onValueChange={(value) => setFilters({ ...filters, categoryId: value || null })}
                disabled={isLoading}
              >
                <SelectTrigger className='mt-2 h-10 w-full'>
                  <SelectValue placeholder='Select category' />
                </SelectTrigger>
                <SelectContent>
                  {categories.map((category) => (
                    <SelectItem
                      key={category.id}
                      value={category.id}
                    >
                      {category.categoryName}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className='w-full space-y-2 sm:flex-1'>
              <label className='text-sm font-medium'>Status</label>
              <Select
                value={filters.statusBill || ''}
                onValueChange={(value) => setFilters({ ...filters, statusBill: value as 'Unpaid' | 'Paid' | null })}
              >
                <SelectTrigger className='mt-2 h-10 w-full'>
                  <SelectValue placeholder='Select status' />
                </SelectTrigger>
                <SelectContent>
                  {statusOptions.map((option) => (
                    <SelectItem
                      key={option.id}
                      value={option.id}
                    >
                      {option.statusName}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className='flex-shrink-0'>
            {hasActiveFilters && (
              <Button
                variant='outline'
                onClick={onReset}
                className='flex h-9 items-center gap-2 px-4'
              >
                <FilterX className='h-4 w-4' />
                Clear Filters
              </Button>
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
