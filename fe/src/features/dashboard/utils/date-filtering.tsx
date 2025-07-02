import type { AssetsType } from '@/features/assets/view-all-assets/model'
import type { DateFilter } from '../model'

export const filterAssetsByDate = (assets: AssetsType[], filter: DateFilter): AssetsType[] => {
  if (!assets || assets.length === 0) return []

  return assets.filter((asset) => {
    let targetDate: Date

    switch (filter.dateField) {
      case 'purchase':
        targetDate = new Date(asset.purchaseDate)
        break
      case 'warranty':
        targetDate = new Date(asset.warrantExpiry)
        break
      case 'created':
        targetDate = new Date(asset.purchaseDate)
        break
      default:
        targetDate = new Date(asset.purchaseDate)
    }

    if (filter.month && filter.year) {
      return targetDate.getMonth() + 1 === filter.month && targetDate.getFullYear() === filter.year
    }

    if (!filter.month && filter.year) {
      return targetDate.getFullYear() === filter.year
    }

    if (filter.month && !filter.year) {
      return targetDate.getMonth() + 1 === filter.month
    }

    if (!filter.month && !filter.year) {
      if (filter.startDate && filter.endDate) {
        const normalizedTargetDate = new Date(targetDate.getFullYear(), targetDate.getMonth(), targetDate.getDate())
        const normalizedStartDate = new Date(
          filter.startDate.getFullYear(),
          filter.startDate.getMonth(),
          filter.startDate.getDate()
        )
        const normalizedEndDate = new Date(
          filter.endDate.getFullYear(),
          filter.endDate.getMonth(),
          filter.endDate.getDate()
        )

        return normalizedTargetDate >= normalizedStartDate && normalizedTargetDate <= normalizedEndDate
      }
      return true
    }

    return true
  })
}

export const getDateRangeText = (filter: DateFilter): string => {
  if (filter.month && filter.year) {
    const monthNames = [
      'January',
      'February',
      'March',
      'April',
      'May',
      'June',
      'July',
      'August',
      'September',
      'October',
      'November',
      'December',
    ]
    return `${monthNames[filter.month - 1]} ${filter.year}`
  }

  if (!filter.month && filter.year) {
    return `All of ${filter.year}`
  }

  if (filter.month && !filter.year) {
    const monthNames = [
      'January',
      'February',
      'March',
      'April',
      'May',
      'June',
      'July',
      'August',
      'September',
      'October',
      'November',
      'December',
    ]
    return `${monthNames[filter.month - 1]} (All Years)`
  }

  if (filter.startDate && filter.endDate) {
    return `${filter.startDate.toLocaleDateString()} - ${filter.endDate.toLocaleDateString()}`
  }

  return 'All Time'
}
export const getLast10Years = (): number[] => {
  const currentYear = new Date().getFullYear()
  const years: number[] = []

  for (let i = 0; i < 10; i++) {
    years.push(currentYear - i)
  }

  return years
}
export const getAvailableYears = (assets: AssetsType[]): number[] => {
  if (!assets || assets.length === 0) return []

  const years = new Set<number>()

  assets.forEach((asset) => {
    const purchaseYear = new Date(asset.purchaseDate).getFullYear()
    const warrantyYear = new Date(asset.warrantExpiry).getFullYear()

    years.add(purchaseYear)
    years.add(warrantyYear)
  })

  return Array.from(years).sort((a, b) => b - a)
}

export const getAvailableMonths = (): Array<{ value: number; label: string }> => {
  return [
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
}

export const getPreviousMonth = (month: number, year: number): { month: number; year: number } => {
  if (month === 1) {
    return { month: 12, year: year - 1 }
  }
  return { month: month - 1, year }
}

export const getPreviousYear = (month: number, year: number): { month: number; year: number } => {
  return { month, year: year - 1 }
}
