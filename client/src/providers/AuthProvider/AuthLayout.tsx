import { Await, useLoaderData, useOutlet } from "react-router-dom";
import { AuthProvider, AuthenticationInfo } from "./AuthProvider";
import { Suspense } from "react";

export const AuthLayout = () => {
  const outlet = useOutlet();
  const { userData } = useLoaderData() as { userData: AuthenticationInfo };

  return (
    <Suspense fallback={<div>Loading...</div>}>
      <Await
        resolve={userData}
        children={(authenticationInfo) => {
          return (
            <AuthProvider authenticationInfo={authenticationInfo}>
              {outlet}
            </AuthProvider>
          );
        }}
      />
    </Suspense>
  );
};
