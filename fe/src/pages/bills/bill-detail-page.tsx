import { useParams, useNavigate } from 'react-router-dom'
import { useEffect, useState } from 'react'
import { BillDetailModal } from '@/features/bills/components/bill-detail-modal'
import { getBillByNumber } from '@/features/bills/api/get-bill-by-number'
import { tryCatch } from '@/utils'
import type { BillType } from '@/features/bills/model/bill-types'
import { toast } from 'sonner'

const BillDetailPage = () => {
  const { billNumber } = useParams()
  const navigate = useNavigate()
  const [bill, setBill] = useState<BillType | null>(null)
  const [open, setOpen] = useState(true)

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

      setBill(response.data.data)
    }

    fetchBill()
  }, [billNumber, navigate])

  const handleClose = () => {
    setOpen(false)
    navigate('/bills')
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
