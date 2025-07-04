import { LayoutDashboard, NotebookPen, Computer, SendToBack, PersonStanding, Receipt } from 'lucide-react'
import { useAppSelector } from '@/hooks'

export const SidebarNav = () => {
  const user = useAppSelector((state) => state.auth.user)
  const role = user.role.slug
  const loading = useAppSelector((state) => state.auth.loading)

  if (loading || !role) {
    return {
      navMain: [],
    }
  }
  const canManageBills = role === 'admin'

  const canCreateAsset = role === 'admin' || (role === 'assetManager' && 'limited scope')
  const canUpdateSchedule = role === 'admin' || role === 'assetManager'
  const canCreateSchedule = role === 'admin' || role === 'assetManager'
  const canTransferRequests = role === 'admin' || role === 'assetManager'
  const canViewAssignments = role === 'admin' || role === 'assetManager'
  const canCreateTransfer = role === 'viewer'
  const canNotViewTransfer = role !== 'viewer'
  const data = {
    navMain: [
      {
        title: 'Dashboard',
        url: '/',
        icon: LayoutDashboard,
        haveChildren: false,
        show: true,
      },
      {
        title: 'Bills',
        url: '/bills',
        icon: Receipt,
        haveChildren: false,
        show: canManageBills,
      },
      {
        title: 'Assets',
        url: '/assets',
        icon: Computer,
        active: true,
        haveChildren: true,
        show: true,
        items: [
          {
            title: 'All assets',
            url: '/assets',
          },
          {
            title: 'Create assets',
            url: '/assets/create-asset',
            show: canCreateAsset,
          },
          {
            title: 'Maintenance schedule',
            url: '/assets/maintenance-schedule',
            show: canUpdateSchedule,
          },
          {
            title: 'Create maintenance schedule',
            url: '/assets/maintenance-schedule/create',
            show: canCreateSchedule,
          },
        ].filter((item) => item.show !== false),
      },

      {
        title: 'Assignments',
        url: '/assignments',
        icon: NotebookPen,
        active: true,
        haveChildren: true,
        show: canViewAssignments,
        items: [
          {
            title: 'All assignments',
            url: '/assignments',
          },
        ],
      },
      {
        title: 'Transfers',
        url: '/transfers',
        icon: SendToBack,
        active: true,
        haveChildren: true,
        show: canTransferRequests,
        items: [
          {
            title: 'All transfers',
            url: '/transfers',
            show: canNotViewTransfer,
          },
          {
            title: 'Create new request transfer',
            url: '/transfers/create-request-transfer',
            show: canCreateTransfer,
          },
        ].filter((item) => item.show !== false),
      },
      {
        title: 'User',
        url: '/user',
        icon: PersonStanding,
        active: true,
        haveChildren: true,
        show: user.role.slug === 'admin',
        items: [
          {
            title: 'Assign role',
            url: '/user/assign-role',
          },
          {
            title: 'Assign department',
            url: '/user/assign-department',
          },
        ],
      },
    ].filter((item) => item.show !== false),
  }
  return data
}
