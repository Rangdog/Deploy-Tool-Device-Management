import { FormProvider, useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { type DataLoginType, loginSchema } from './model/schema'
import { Form, FormButtonSubmit, FormInput } from '@/components/ui'
import { useNavigate, useLocation } from 'react-router-dom'
import { toast } from 'sonner'
import { useAppDispatch } from '@/hooks'
import { logIn } from '../slice'
import { useState } from 'react'
import Cookies from 'js-cookie'
import { LinkToSignUp, RememberMeField, LinkToForgetPassword } from './_components'
import { Lock, Mail } from 'lucide-react'

const LoginForm = () => {
  const [isPending, setIsPending] = useState(false)
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const form = useForm<DataLoginType>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: '',
      password: '',
      rememberMe: false,
    },
  })
  const location = useLocation()
  const params = new URLSearchParams(location.search)
  const redirectPath = params.get('redirect')

  const onSubmit = async (data: DataLoginType) => {
    setIsPending(true)
    const result = await dispatch(logIn(data)).unwrap()

    if (!result.success) {
      toast.error(result.error?.message || 'Login failed, please try again')
      return
    }
    if (result.data.is_active === false) {
      toast.error('User is inactive, please go to email and verify your account')
      return
    }
    Cookies.set('accessToken', result.data.access_token)
    Cookies.set('refreshToken', result.data.refresh_token)
    toast.success('Login successfully')
    navigate(redirectPath || '/', { replace: true })
    setIsPending(false)
  }
  return (
    <div>
      <div className='mb-4 flex flex-col items-center sm:mb-6'>
        <h1 className='text-foreground mb-1 text-xl font-bold sm:mb-2 sm:text-2xl'>Sign in to your account</h1>
        <p className='text-muted-foreground px-2 text-center text-sm'>Manage your devices and assets securely</p>
      </div>

      <FormProvider {...form}>
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(onSubmit)}
            className='space-y-4 sm:space-y-5'
          >
            <FormInput
              name='email'
              type='email'
              label='Email Address'
              Icon={Mail}
              placeholder='Enter your email address'
            />

            <FormInput
              name='password'
              type='password'
              label='Password'
              Icon={Lock}
              placeholder='••••••••'
            />
            <RememberMeField form={form} />
            <FormButtonSubmit
              className='bg-primary hover:bg-primary/90 mt-5 flex h-9 w-full items-center justify-center gap-1.5 text-sm font-medium sm:mt-6 sm:h-10 sm:gap-2 sm:text-base'
              isPending={isPending}
              type='Log In'
              onSubmit={onSubmit}
            />
            <div className='flex flex-col items-center gap-2 sm:flex-row sm:justify-between'>
              <LinkToSignUp />
              <LinkToForgetPassword />
            </div>
          </form>
        </Form>
      </FormProvider>
    </div>
  )
}

export default LoginForm
