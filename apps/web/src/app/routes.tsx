import { Navigate, Outlet, RouteObject } from "react-router-dom";
import { FeedPage } from "../features/feed/pages/FeedPage";
import { LoginPage } from "../features/auth/pages/LoginPage";
import { SignupPage } from "../features/auth/pages/SignupPage";
import { MyProfilePage } from "../features/profile/pages/MyProfilePage";
import { PublicProfilePage } from "../features/profile/pages/PublicProfilePage";
import { useSession } from "../features/auth/hooks/useSession";

function RequireAuth() {
  const { user, isHydrating } = useSession();

  if (isHydrating) {
    return <p style={{ padding: "1rem" }}>Loading session...</p>;
  }

  return user ? <Outlet /> : <Navigate to="/login" replace />;
}

function GuestOnly() {
  const { user, isHydrating } = useSession();

  if (isHydrating) {
    return <p style={{ padding: "1rem" }}>Loading session...</p>;
  }

  return user ? <Navigate to="/feed" replace /> : <Outlet />;
}

export const appRoutes: RouteObject[] = [
  { path: "/", element: <Navigate to="/feed" replace /> },
  {
    element: <GuestOnly />,
    children: [
      { path: "login", element: <LoginPage /> },
      { path: "signup", element: <SignupPage /> },
    ],
  },
  {
    element: <RequireAuth />,
    children: [
      { path: "feed", element: <FeedPage /> },
      { path: "profile/me", element: <MyProfilePage /> },
      { path: "profile/:userId", element: <PublicProfilePage /> },
    ],
  },
];
