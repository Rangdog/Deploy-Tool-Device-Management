import { Button } from '@/components'
import { Plus } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { useAppSelector } from '@/hooks'
export const ButtonCreateNewAssets = () => {
  const navigate = useNavigate()
  const user = useAppSelector((state) => state.auth.user)

  if (user.role.slug === 'viewer') {
    return null
  }
  return (
    <Button
      size='sm'
      onClick={() => navigate('/assets/create-asset')}
    >
      <Plus className='mr-2 h-4 w-4' />
      Create New Asset
    </Button>
  )
}
