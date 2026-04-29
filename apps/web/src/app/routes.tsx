import React from "react";
import { PropsWithChildren } from "react";
import { Navigate, RouteObject } from "react-router-dom";
import { FeedPage } from "../features/feed/pages/FeedPage";
import { LoginPage } from "../features/auth/pages/LoginPage";
import { SignupPage } from "../features/auth/pages/SignupPage";
import { MyProfilePage } from "../features/profile/pages/MyProfilePage";
import { PublicProfilePage } from "../features/profile/pages/PublicProfilePage";
import { useSession } from "../features/auth/hooks/useSession";

function RequireAuth({ children }: PropsWithChildren) {
  const { user, isHydrating } = useSession();

  if (isHydrating) {
    return <p style={{ padding: "1rem" }}>Loading session...</p>;
  }

  return user ? <>{children}</> : <Navigate to="/login" replace />;
}

function GuestOnly({ children }: PropsWithChildren) {
  const { user, isHydrating } = useSession();

  if (isHydrating) {
    return <p style={{ padding: "1rem" }}>Loading session...</p>;
  }

  return user ? <Navigate to="/feed" replace /> : <>{children}</>;
}

export const appRoutes: RouteObject[] = [
  { path: "/", element: <Navigate to="/feed" replace /> },
  { path: "/login", element: <GuestOnly><LoginPage /></GuestOnly> },
  { path: "/signup", element: <GuestOnly><SignupPage /></GuestOnly> },
  { path: "/feed", element: <RequireAuth><FeedPage /></RequireAuth> },
  { path: "/profile/me", element: <RequireAuth><MyProfilePage /></RequireAuth> },
  { path: "/profile/:userId", element: <RequireAuth><PublicProfilePage /></RequireAuth> },
  { path: "*", element: <Navigate to="/login" replace /> },
];
