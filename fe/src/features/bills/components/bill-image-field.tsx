import { FormControl, FormField, FormItem, FormMessage, FormLabel, Input } from '@/components/ui'
import { ImageIcon } from 'lucide-react'
import type { CreateBillFormType } from '../model/create-bill-schema'
import type { UseFormReturn } from 'react-hook-form'

export const BillImageField = ({
  form,
  imageName,
  setImageName,
}: {
  form: UseFormReturn<CreateBillFormType>
  imageName: string
  setImageName: React.Dispatch<React.SetStateAction<string>>
}) => {
  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      form.setValue('imageUpload', file)
      setImageName(file.name)
    }
  }

  return (
    <FormField
      control={form.control}
      name='imageUpload'
      render={() => (
        <FormItem>
          <FormLabel>Image Upload</FormLabel>
          <FormControl>
            <div className='flex w-full flex-col items-center justify-center'>
              <label
                htmlFor='bill-image-upload'
                className='flex h-28 w-full cursor-pointer flex-col items-center justify-center overflow-hidden rounded-lg border-2 border-dashed bg-gray-50 hover:bg-gray-100 dark:border-gray-600 dark:bg-slate-900 dark:hover:border-gray-500'
              >
                <div className='flex flex-col items-center justify-center pt-5 pb-6'>
                  <ImageIcon className='mb-2 h-8 w-8 text-gray-500 dark:text-gray-400' />
                  <p className='mb-1 text-sm text-gray-500 dark:text-gray-400'>
                    {imageName ? imageName : 'Click to upload image'}
                  </p>
                  <p className='text-xs text-gray-500 dark:text-gray-400'>PNG, JPG, GIF up to 10MB</p>
                </div>
                <Input
                  id='bill-image-upload'
                  type='file'
                  accept='image/*'
                  className='hidden'
                  onChange={handleImageChange}
                />
              </label>
            </div>
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  )
}
