# Tool Device Management - Frontend

A modern web application for managing assets, users, departments, and maintenance in your organization. Built with React, TypeScript, and Vite for a fast, scalable, and maintainable experience.

## 🌟 Features

- **Dashboard**: Overview of asset statistics, recent purchased devices, and quick insights.
- **Asset Management**: View, create, update, and filter assets with categories, departments, and statuses.
- **Assignment Management**: Assign assets to users and track assignment history.
- **Transfer Requests**: Create and manage asset transfer requests between departments.
- **User Management**: Assign roles and departments to users.
- **Department Management**: Organize users and assets by department.
- **Maintenance Scheduling**: Plan and track maintenance schedules for assets.
- **Notifications**: Receive realtime update notifications on assignments, transfers, and maintenance.
- **Authentication & Authorization**: Secure login and role-based access control.

## 🚀 Tech Stack

- [React] + [TypeScript]
- [Vite] for fast development
- [Redux Toolkit] [Redux Thunk] for state management
- [Zod] [React Hook Form] for form and validation
- [Tailwind CSS] for styling
- [Shadcn] for UI components

## 📦 Getting Started

1. **Install dependencies:**
   ```bash
   npm install
   ```
2. **Start the development server:**
   ```bash
   npm run dev
   ```
3. **Open your browser:**
   Visit [http://localhost:5173](http://localhost:5173) to view the app.

## 🌐 Production

You can access the live application here:

[Production Website](https://tool-device-management.vercel.app/)

## 🧑‍💻 For Developers

- Code is organized by feature in `src/features/` for scalability.
- Reusable UI components are located in `src/components/ui/`.
- API calls and models are separated for clarity.
- ESLint and Prettier are configured for code quality.

## 📁 Project Structure

```
fe/
  src/
    features/        # Feature modules (assets, user, assignments, ...)
    components/      # Reusable UI components
    hooks/           # Custom React hooks
    lib/             # Utilities and HTTP client
    pages/           # Page components
    redux-store/     # Redux store setup
    router/          # App routing
    utils/           # Helper functions
  public/            # Static assets
  ...
```

## 📜 Future Improvements

`Dark mode and theme customization support.` (Done)
`User can update their personal information like upload avatar, change information. ` (Done)
`Enable side-by-side comparison of multiple assets, including specs, costs, and performance ` (Done)
`Admins can view statistical reports filtered by month and year. ` (Done)
`Admins can create bills after creating assets, and generate a monthly summary bill for printing at the end of the month. ` (Done)
`Admin (save records when buy, update,... ) ` (Done)
`Manager (assign asset to employee) ` (Done)
`Employee (see assets ,` `assign asset ,` (Done)
`Admin doesn't have to approve ` (Done)
**Admin assign asset -> Manager -> Employee** (In-progress)
**Employee assign to another employee in the same department** (In-progress)
**Improve UI/UX: Form create , update asset, button, font, size,....** (In-progress)
