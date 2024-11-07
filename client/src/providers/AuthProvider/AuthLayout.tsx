import { Suspense } from 'react';
import { Await, useLoaderData, useOutlet } from 'react-router-dom';

import { AuthProvider, AuthenticationInfo } from './AuthProvider';
import Spinner from 'components/Spinner/Spinner';

export const AuthLayout = () => {
  const outlet = useOutlet();
  const { authData } = useLoaderData() as {
    authData: AuthenticationInfo | null;
  };

  return (
    <Suspense fallback={<Spinner />}>
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
