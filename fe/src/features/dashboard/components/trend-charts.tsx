import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui'
import { ReusablePieChart } from '@/components/ui/charts/pie-chart'
import {
  LineChart,
  Line,
  AreaChart,
  Area,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts'
import type { TrendData, StatusStats, DepartmentStats } from '../model'
import { formatCurrency, formatNumber } from '../utils'
import { ChartBarIcon, PieChart } from 'lucide-react'

interface TrendChartsProps {
  trendData: TrendData[]
  departmentData: DepartmentStats[]
  statusData: StatusStats[]
  className?: string
  isPending?: boolean
  showOverview?: boolean
}

const CustomTooltip = ({ active, payload, label }: any) => {
  if (active && payload && payload.length) {
    return (
      <div className='rounded-lg border bg-white p-3 shadow-lg dark:border-gray-700 dark:bg-gray-800'>
        <p className='font-medium dark:text-gray-200'>{label}</p>
        {payload.map((entry: any, index: number) => (
          <p
            key={index}
            style={{ color: entry.color }}
            className='dark:text-opacity-90'
          >
            {entry.name}:{' '}
            {entry.name.includes('Value') || entry.name.includes('Cost')
              ? formatCurrency(entry.value)
              : formatNumber(entry.value)}
          </p>
        ))}
      </div>
    )
  }
  return null
}

export const TrendCharts = ({
  trendData,
  departmentData,
  statusData,
  className = '',
  isPending = false,
}: TrendChartsProps) => {
  const transformedStatusData = statusData.map((item) => {
    let name = item.status

    if (item.status === 'Retired' || item.status === 'Disposed') {
      name = 'Retired / Disposed'
    }
    return { name, value: item.count }
  })

  const groupedData = transformedStatusData.reduce(
    (acc, curr) => {
      const existing = acc.find((item) => item.name === curr.name)
      if (existing) {
        existing.value += curr.value
      } else {
        acc.push(curr)
      }
      return acc
    },
    [] as Array<{ name: string; value: number }>
  )

  return (
    <div className={`space-y-4 sm:space-y-6 ${className}`}>
      <div className='grid grid-cols-1 gap-4 lg:grid-cols-2'>
        <Card className='w-full transition-all duration-200 dark:border-gray-700 dark:bg-gray-800/50'>
          <CardHeader className='p-3 sm:p-4'>
            <CardTitle className='text-base sm:text-lg dark:text-gray-200'>
              <div className='flex items-center gap-2'>
                <ChartBarIcon className='h-4 w-4 sm:h-5 sm:w-5' />
                Asset Count Trend
              </div>
            </CardTitle>
          </CardHeader>
          <CardContent className='p-3 pt-0 sm:p-4'>
            <div className='h-[200px] sm:h-[250px] md:h-[300px]'>
              <ResponsiveContainer
                width='100%'
                height='100%'
              >
                <LineChart data={trendData}>
                  <CartesianGrid
                    strokeDasharray='3 3'
                    stroke='#e2e8f0'
                    className='dark:opacity-30'
                  />
                  <XAxis
                    dataKey='month'
                    tick={{ fontSize: 12, dy: 8 }}
                    stroke='#64748b'
                    className='dark:text-gray-400'
                    height={40}
                  />
                  <YAxis
                    tick={{ fontSize: 12 }}
                    stroke='#64748b'
                    className='dark:text-gray-400'
                    width={40}
                  />
                  <Tooltip content={<CustomTooltip />} />
                  <Line
                    type='monotone'
                    dataKey='totalAssets'
                    stroke='#3b82f6'
                    strokeWidth={2}
                    name='Total Assets'
                    dot={{ fill: '#3b82f6', strokeWidth: 2, r: 3 }}
                  />
                  <Line
                    type='monotone'
                    dataKey='newAssets'
                    stroke='#22c55e'
                    strokeWidth={2}
                    name='New Assets'
                    dot={{ fill: '#22c55e', strokeWidth: 2, r: 3 }}
                  />
                </LineChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>

        <Card className='w-full transition-all duration-200 dark:border-gray-700 dark:bg-gray-800/50'>
          <CardHeader className='p-3 sm:p-4'>
            <CardTitle className='text-base sm:text-lg dark:text-gray-200'>
              <div className='flex items-center gap-2'>
                <ChartBarIcon className='h-4 w-4 sm:h-5 sm:w-5' />
                Asset Value Trend
              </div>
            </CardTitle>
          </CardHeader>
          <CardContent className='p-3 pt-0 sm:p-4'>
            <div className='h-[200px] sm:h-[250px] md:h-[300px]'>
              <ResponsiveContainer
                width='100%'
                height='100%'
              >
                <AreaChart data={trendData}>
                  <CartesianGrid
                    strokeDasharray='3 3'
                    stroke='#e2e8f0'
                    className='dark:opacity-30'
                  />
                  <XAxis
                    dataKey='month'
                    tick={{ fontSize: 12, dy: 8 }}
                    stroke='#64748b'
                    className='dark:text-gray-400'
                    height={40}
                  />
                  <YAxis
                    tick={{ fontSize: 12 }}
                    tickFormatter={(value) => formatCurrency(value)}
                    stroke='#64748b'
                    className='dark:text-gray-400'
                    width={60}
                  />
                  <Tooltip content={<CustomTooltip />} />
                  <Area
                    type='monotone'
                    dataKey='totalValue'
                    stroke='#8b5cf6'
                    fill='#8b5cf6'
                    fillOpacity={0.3}
                    name='Total Value'
                  />
                </AreaChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>
      </div>

      <div className='grid grid-cols-1 gap-4 sm:gap-6 lg:grid-cols-2'>
        <Card className='transition-all duration-200 dark:border-gray-700 dark:bg-gray-800/50'>
          <CardHeader>
            <CardTitle className='text-base sm:text-lg dark:text-gray-200'>
              <div className='flex items-center gap-2'>
                <ChartBarIcon className='h-5 w-5' />
                Assets by Department
              </div>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className='h-[250px] sm:h-[300px]'>
              <ResponsiveContainer
                width='100%'
                height='100%'
              >
                <BarChart
                  data={departmentData}
                  margin={{ right: 20 }}
                >
                  <CartesianGrid
                    strokeDasharray='3 3'
                    stroke='#e2e8f0'
                    className='dark:opacity-30'
                  />
                  <XAxis
                    dataKey='departmentName'
                    tick={{ fontSize: 12 }}
                    angle={-45}
                    textAnchor='end'
                    height={30}
                    stroke='#64748b'
                    className='dark:text-gray-400'
                  />
                  <YAxis
                    tick={{ fontSize: 12 }}
                    stroke='#64748b'
                    className='dark:text-gray-400'
                  />
                  <Tooltip content={<CustomTooltip />} />
                  <Bar
                    dataKey='count'
                    fill='#3b82f6'
                    name='Asset Count'
                    radius={[4, 4, 0, 0]}
                  />
                </BarChart>
              </ResponsiveContainer>
            </div>
            <div className='mt-4 ml-10 text-left text-xs text-gray-500 dark:text-gray-400'>
              <p>* Note: Department only displays data by date</p>
            </div>
          </CardContent>
        </Card>

        <ReusablePieChart
          icon={<PieChart className='h-5 w-5' />}
          data={groupedData}
          title='Asset Status Distribution'
          isPending={isPending}
          showAnimation={true}
          animationDelay={2000}
          height='h-[250px] sm:h-[300px]'
          className='dark:border-gray-700 dark:bg-gray-800/50'
        />
      </div>
    </div>
  )
}
