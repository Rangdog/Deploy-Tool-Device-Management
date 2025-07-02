import { useState } from 'react'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Button,
  Popover,
  PopoverTrigger,
  PopoverContent,
  Calendar,
} from '@/components/ui'
import { CalendarIcon, X, Calendar as CalendarDays, Clock } from 'lucide-react'
import { format } from 'date-fns'
import type { DateFilter as DateFilterType } from '../model/statistics-types'
import { getLast10Years } from '../utils'
import type { AssetsType } from '@/features/assets/view-all-assets/model'

interface DateFilterProps {
  dateFilter: DateFilterType
  onDateFilterChange: (filter: DateFilterType) => void
  originalAssets: AssetsType[]
  assets: AssetsType[]
  className?: string
  showClearButton?: boolean
}

export const DateFilter = ({
  dateFilter,
  onDateFilterChange,
  // originalAssets,
  // assets,
  className = '',
  showClearButton = false,
}: DateFilterProps) => {
  const getCurrentMode = (): 'month-year' | 'date-range' => {
    if (dateFilter.startDate && dateFilter.endDate) return 'date-range'
    return 'month-year'
  }

  const [filterMode, setFilterMode] = useState<'month-year' | 'date-range'>(getCurrentMode())

  const availableYears = getLast10Years()
  const availableMonths = [
    { value: 1, label: 'January' },
    { value: 2, label: 'February' },
    { value: 3, label: 'March' },
    { value: 4, label: 'April' },
    { value: 5, label: 'May' },
    { value: 6, label: 'June' },
    { value: 7, label: 'July' },
    { value: 8, label: 'August' },
    { value: 9, label: 'September' },
    { value: 10, label: 'October' },
    { value: 11, label: 'November' },
    { value: 12, label: 'December' },
  ]

  const handleModeSwitch = (mode: 'month-year' | 'date-range') => {
    setFilterMode(mode)

    if (mode === 'month-year') {
      onDateFilterChange({
        ...dateFilter,
        startDate: undefined,
        endDate: undefined,
        singleDate: undefined,
      })
    } else {
      onDateFilterChange({
        ...dateFilter,
        month: undefined,
        year: undefined,
        singleDate: undefined,
      })
    }
  }

  const handleClearFilter = () => {
    if (filterMode === 'month-year') {
      onDateFilterChange({
        ...dateFilter,
        month: undefined,
        year: undefined,
      })
    } else {
      onDateFilterChange({
        ...dateFilter,
        startDate: undefined,
        endDate: undefined,
      })
    }
  }

  const hasActiveFilters =
    filterMode === 'month-year'
      ? !!(dateFilter.month || dateFilter.year)
      : !!(dateFilter.startDate && dateFilter.endDate)

  return (
    <div className={`flex flex-wrap items-end gap-4 ${className}`}>
      <div className='flex flex-col gap-2'>
        <label className='text-sm font-medium text-black dark:text-white'>Filter Type</label>
        <div className='bg-background flex h-9 rounded-md border'>
          <Button
            variant={filterMode === 'month-year' ? 'default' : 'ghost'}
            size='sm'
            onClick={() => handleModeSwitch('month-year')}
            className='h-full flex-1 rounded-r-none text-xs'
          >
            <Clock className='mr-1 h-3 w-3' />
            Month/Year
          </Button>
          <Button
            variant={filterMode === 'date-range' ? 'default' : 'ghost'}
            size='sm'
            onClick={() => handleModeSwitch('date-range')}
            className='h-full flex-1 rounded-l-none text-xs'
          >
            <CalendarDays className='mr-1 h-3 w-3' />
            Date Range
          </Button>
        </div>
      </div>

      {filterMode === 'month-year' && (
        <>
          <div className='flex flex-col gap-2'>
            <label className='text-sm font-medium text-black dark:text-white'>Year</label>
            <Select
              value={dateFilter.year?.toString() || 'all'}
              onValueChange={(value) =>
                onDateFilterChange({
                  ...dateFilter,
                  year: value === 'all' ? undefined : parseInt(value),
                })
              }
            >
              <SelectTrigger className='h-10 w-30 sm:w-45'>
                <SelectValue placeholder='All Years' />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value='all'>All Years</SelectItem>
                {availableYears.map((year) => (
                  <SelectItem
                    key={year}
                    value={year.toString()}
                  >
                    {year}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className='flex flex-col gap-2'>
            <label className='text-sm font-medium text-black dark:text-white'>Month</label>
            <Select
              value={dateFilter.month?.toString() || 'all'}
              onValueChange={(value) =>
                onDateFilterChange({
                  ...dateFilter,
                  month: value === 'all' ? undefined : parseInt(value),
                })
              }
            >
              <SelectTrigger className='h-10 w-30 sm:w-45'>
                <SelectValue placeholder='All Months' />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value='all'>All Months</SelectItem>
                {availableMonths.map((month) => (
                  <SelectItem
                    key={month.value}
                    value={month.value.toString()}
                  >
                    {month.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </>
      )}

      {filterMode === 'date-range' && (
        <>
          <div className='flex flex-col gap-2'>
            <label className='text-sm font-medium text-black dark:text-white'>Start Date</label>
            <Popover>
              <PopoverTrigger asChild>
                <Button
                  variant='outline'
                  className={`h-9 w-30 justify-start text-left font-normal sm:w-45 ${
                    !dateFilter.startDate
                      ? 'text-muted-foreground dark:text-muted-foreground'
                      : 'text-black dark:text-white'
                  }`}
                >
                  <CalendarIcon className='mr-1 h-4 w-4' />
                  {dateFilter.startDate ? format(dateFilter.startDate, 'MMM d, yyyy') : 'Start date'}
                </Button>
              </PopoverTrigger>
              <PopoverContent className='w-auto p-0'>
                <Calendar
                  mode='single'
                  selected={dateFilter.startDate}
                  onSelect={(date) =>
                    onDateFilterChange({
                      ...dateFilter,
                      startDate: date,
                    })
                  }
                  initialFocus
                />
              </PopoverContent>
            </Popover>
          </div>

          <div className='flex flex-col gap-2'>
            <label className='text-sm font-medium text-black dark:text-white'>End Date</label>
            <Popover>
              <PopoverTrigger asChild>
                <Button
                  variant='outline'
                  className={`h-9 w-30 justify-start text-left font-normal sm:w-45 ${
                    !dateFilter.endDate
                      ? 'text-muted-foreground dark:text-muted-foreground'
                      : 'text-black dark:text-white'
                  }`}
                >
                  <CalendarIcon className='mr-1 h-4 w-4' />
                  {dateFilter.endDate ? format(dateFilter.endDate, 'MMM d, yyyy') : 'End date'}
                </Button>
              </PopoverTrigger>
              <PopoverContent className='w-auto p-0'>
                <Calendar
                  mode='single'
                  selected={dateFilter.endDate}
                  onSelect={(date) =>
                    onDateFilterChange({
                      ...dateFilter,
                      endDate: date,
                    })
                  }
                  initialFocus
                />
              </PopoverContent>
            </Popover>
          </div>
        </>
      )}

      {showClearButton && hasActiveFilters && (
        <Button
          variant='outline'
          size='sm'
          onClick={handleClearFilter}
          className='h-10'
        >
          <X className='mr-1 h-4 w-4' />
          Clear
        </Button>
      )}
    </div>
  )
}
