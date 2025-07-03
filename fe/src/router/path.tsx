import { createBrowserRouter, Navigate } from 'react-router-dom'
import {
  DashboardPage,
  LoginPage,
  SignupPage,
  ForgetPasswordPage,
  ResetPasswordPage,
  AssetsPage,
  CreateAssetPage,
  AssetDetailPage,
  AssignmentsPage,
  UpdateAssetPage,
  AssignmentDetailPage,
  AllRequestTransferPage,
  RequestTransferDetailPage,
  CreateNewRequestTransferPage,
  ViewAllMaintenanceSchedulePage,
  CreateAssetMaintenanceSchedulePage,
  AssetComparisonPage,
  StatisticReportPage,
} from '@/pages/dashboard'
import { ProtectedRoute, AuthRoute } from './auth-guard'
import { BillsManagement } from '@/features/bills'
import Layout from '@/layout/layout'
import { AssignDepartmentForUserPage, AssignRoleForUserPage, EditProfilePage } from '@/pages/user'

export const router = createBrowserRouter([
  {
    path: '/',
    element: <ProtectedRoute />,
    children: [
      {
        path: '/',
        element: <Layout />,
        children: [
          {
            path: 'dashboard',
            element: <DashboardPage />,
          },
          {
            path: 'statistics-reports',
            element: <StatisticReportPage />,
          },
          {
            path: '',
            element: <Navigate to='/dashboard' />,
          },
          {
            path: 'assets',
            children: [
              {
                path: '',
                element: <AssetsPage />,
              },
              {
                path: 'create-asset',
                element: <CreateAssetPage />,
              },
              {
                path: ':id',
                element: <AssetDetailPage />,
              },
              {
                path: 'update/:id',
                element: <UpdateAssetPage />,
              },
              {
                path: 'asset-comparison',
                element: <AssetComparisonPage />,
              },
              {
                path: 'maintenance-schedule',
                children: [
                  {
                    path: '',
                    element: <ViewAllMaintenanceSchedulePage />,
                  },
                  {
                    path: 'create',
                    element: <CreateAssetMaintenanceSchedulePage />,
                  },
                ],
              },
            ],
          },
          {
            path: 'assignments',
            children: [
              {
                path: '',
                element: <AssignmentsPage />,
              },
              {
                path: ':id',
                element: <AssignmentDetailPage />,
              },
            ],
          },
          {
            path: 'transfers',
            children: [
              {
                path: '',
                element: <AllRequestTransferPage />,
              },
              {
                path: ':id',
                element: <RequestTransferDetailPage />,
              },
              {
                path: 'create-request-transfer',
                element: <CreateNewRequestTransferPage />,
              },
            ],
          },
          {
            path: 'user',
            children: [
              {
                path: 'assign-role',
                element: <AssignRoleForUserPage />,
              },
              {
                path: 'assign-department',
                element: <AssignDepartmentForUserPage />,
              },
              {
                path: 'edit-profile',
                element: <EditProfilePage />,
              },
            ],
          },
          {
            path: 'bills',
            element: <BillsManagement />,
          },
        ],
      },
    ],
  },
  {
    path: '',
    element: <AuthRoute />,
    children: [
      {
        path: 'login',
        element: <LoginPage />,
      },
      {
        path: 'signup',
        element: <SignupPage />,
      },
      {
        path: 'forget-password',
        element: <ForgetPasswordPage />,
      },
      {
        path: 'reset-password',
        element: <ResetPasswordPage />,
      },
    ],
  },
])
export const AppPaths = {
  DASHBOARD: '/dashboard',
  STATISTICS_REPORTS: '/statistics-reports',
  ASSETS: '/assets',
  BILLS: '/bills',
  ASSET_COMPARISON: '/assets/asset-comparison',
  CREATE_ASSET: '/assets/create-asset',
  ASSIGNMENTS: '/assignments',
  TRANSFERS: '/transfers',
  LOGIN: '/login',
  SIGNUP: '/signup',
} as const
