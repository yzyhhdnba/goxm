import { AdminReviewPage } from '@/features/admin/pages/AdminReviewPage';
import { Navigate, Route, Routes } from 'react-router-dom';
import { DashboardPage } from '@/features/account/pages/DashboardPage';
import { LoginPage } from '@/features/auth/pages/LoginPage';
import { CreatorStudioPage } from '@/features/creator/pages/CreatorStudioPage';
import { HomePage } from '@/features/feed/pages/HomePage';
import { HistoryPage } from '@/features/history/pages/HistoryPage';
import { NoticesPage } from '@/features/notice/pages/NoticesPage';
import { SearchPage } from '@/features/search/pages/SearchPage';
import { UserProfilePage } from '@/features/user/pages/UserProfilePage';
import { VideoDetailPage } from '@/features/video/pages/VideoDetailPage';
import { AdminRoute, GuestOnlyRoute, ProtectedRoute } from '@/routes/guards';

export function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/search" element={<SearchPage />} />
      <Route path="/sevideo" element={<SearchPage />} />
      <Route path="/seuser" element={<SearchPage />} />
      <Route
        path="/login"
        element={
          <GuestOnlyRoute>
            <LoginPage />
          </GuestOnlyRoute>
        }
      />
      <Route path="/video/:id" element={<VideoDetailPage />} />
      <Route path="/users/:id" element={<UserProfilePage />} />
      <Route
        path="/me"
        element={
          <ProtectedRoute>
            <DashboardPage />
          </ProtectedRoute>
        }
      />
      <Route
        path="/history"
        element={
          <ProtectedRoute>
            <HistoryPage />
          </ProtectedRoute>
        }
      />
      <Route
        path="/creator"
        element={
          <ProtectedRoute>
            <CreatorStudioPage />
          </ProtectedRoute>
        }
      />
      <Route
        path="/notices"
        element={
          <ProtectedRoute>
            <NoticesPage />
          </ProtectedRoute>
        }
      />
      <Route
        path="/admin"
        element={
          <AdminRoute>
            <AdminReviewPage />
          </AdminRoute>
        }
      />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
