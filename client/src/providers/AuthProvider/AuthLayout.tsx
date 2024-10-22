import { Suspense } from 'react';
import { Await, useLoaderData, useOutlet } from 'react-router-dom';

import { AuthProvider, AuthenticationInfo } from './AuthProvider';

export const AuthLayout = () => {
  const outlet = useOutlet();
  const { authData } = useLoaderData() as {
    authData: AuthenticationInfo | null;
  };

  return (
    <Suspense fallback={<div>Loading...</div>}>
      <Await resolve={authData}>
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
