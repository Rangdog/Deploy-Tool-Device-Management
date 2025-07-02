import { Moon, Sun, Palette } from 'lucide-react'
import { useAppDispatch, useAppSelector } from '@/hooks'
import { toggleTheme } from '../slice/theme-slice'
import { SidebarMenuButton, useSidebar } from '@/components/ui/sidebar'

export const DarkModeToggle = () => {
  const dispatch = useAppDispatch()
  const theme = useAppSelector((state) => state.theme.mode)
  const { state, isMobile } = useSidebar()
  const isDark = theme === 'dark'
  const isCollapsed = state === 'collapsed'

  const handleToggle = () => {
    dispatch(toggleTheme())
  }

  const showText = isMobile || !isCollapsed

  return (
    <SidebarMenuButton
      onClick={handleToggle}
      className='w-full justify-between'
      tooltip={!showText ? 'Toggle dark mode' : undefined}
    >
      <div className='flex items-center gap-2'>
        <Palette className='h-4 w-4' />
        {showText && <span className='text-sm'>Dark mode</span>}
      </div>
      {showText && (
        <div
          className={`relative h-5 w-9 rounded-full transition-colors duration-100 ${
            isDark ? 'bg-primary' : 'bg-gray-300'
          }`}
        >
          <div
            className={`bg-background absolute top-0.5 h-4 w-4 rounded-full shadow-md transition-transform duration-100 ${
              isDark ? 'translate-x-4' : 'translate-x-0.5'
            }`}
          >
            <div className='flex h-full w-full items-center justify-center'>
              {isDark ? (
                <Moon className='text-primary h-2.5 w-2.5' />
              ) : (
                <Sun className='text-muted-foreground h-2.5 w-2.5' />
              )}
            </div>
          </div>
        </div>
      )}
    </SidebarMenuButton>
  )
}
