import { Suspense } from 'react';
import { Await, useLoaderData, useOutlet } from 'react-router-dom';

import { AuthProvider, AuthenticationInfo } from './AuthProvider';

export const AuthLayout = () => {
  const outlet = useOutlet();
  const { userData } = useLoaderData() as { userData: AuthenticationInfo };

  return (
    <Suspense fallback={<div>Loading...</div>}>
      <Await resolve={userData}>
        {(authenticationInfo) => {
          return (
            <AuthProvider authenticationInfo={authenticationInfo}>
              {outlet}
            </AuthProvider>
          );
        }}
      </Await>
    </Suspense>
  );
};
