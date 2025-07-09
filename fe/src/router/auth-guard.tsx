import { Navigate, Outlet, useLocation, matchPath } from 'react-router-dom'
import Cookies from 'js-cookie'
import { useGetSession } from '@/hooks/use-get-session'
import { Skeleton } from '@/components'

export const ProtectedRoute = () => {
  const token = Cookies.get('accessToken')
  const { loading } = useGetSession()
  const location = useLocation()

  const isBillDetail = matchPath('/bills/:billNumber', location.pathname)

  if (loading) return <Skeleton className='size-full' />
  if (token) return <Outlet />
  return isBillDetail ? (
    <Navigate
      to={`/login?redirect=${encodeURIComponent(location.pathname + location.search)}`}
      replace
    />
  ) : (
    <Navigate to={'/login'} />
  )
}

export const AuthRoute = () => {
  const token = Cookies.get('accessToken')
  return token ? <Navigate to={'/'} /> : <Outlet />
}
