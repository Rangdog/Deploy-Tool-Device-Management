import { useParams, useNavigate } from 'react-router-dom'
import { useEffect, useState } from 'react'
import { BillDetailModal } from '@/features/bills/components/bill-detail-modal'
import { getBillByNumber } from '@/features/bills/api/get-bill-by-number'
import { tryCatch } from '@/utils'
import type { BillType } from '@/features/bills/model/bill-types'
import { toast } from 'sonner'
import Cookies from 'js-cookie'
import { useAppSelector } from '@/hooks'

const BillDetailPage = () => {
  const { billNumber } = useParams()
  const navigate = useNavigate()
  const [bill, setBill] = useState<BillType | null>(null)
  const [open, setOpen] = useState(true)
  const currentUser = useAppSelector((state) => state.auth.user)

  const getCurrentUserInfo = () => {
    const firstName = Cookies.get('firstName') || currentUser.firstName
    const lastName = Cookies.get('lastName') || currentUser.lastName
    const email = Cookies.get('email') || currentUser.email
    const avatar = Cookies.get('avatar') || currentUser.avatar
    const id = Cookies.get('id') || currentUser.id

    return {
      id: parseInt(id) || parseInt(currentUser.id),
      fullName: `${firstName} ${lastName}`.trim() || 'Current User',
      email: email || 'No email',
      avatar: avatar || null,
    }
  }

  useEffect(() => {
    if (!billNumber) {
      toast.error('Bill number is missing')
      navigate('/bills')
      return
    }

    const fetchBill = async () => {
      const response = await tryCatch(getBillByNumber(billNumber))

      if (response.error || !response.data?.data) {
        toast.error('Failed to fetch bill')
        navigate('/bills')
        return
      }

      const billData = response.data.data
      setBill({
        ...billData,
        creator: billData.creator || getCurrentUserInfo(),
      })
    }

    fetchBill()
  }, [billNumber, navigate])

  const handleClose = () => {
    setOpen(false)
    navigate(-1)
  }

  const handleStatusChange = (billId: number, newStatus: 'Unpaid' | 'Paid') => {
    toast.success(`Bill #${billId} marked as ${newStatus}`)
  }

  return (
    <BillDetailModal
      bill={bill}
      open={open}
      onClose={handleClose}
      onStatusChange={handleStatusChange}
    />
  )
}

export default BillDetailPage
