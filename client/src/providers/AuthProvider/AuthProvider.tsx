import { ManagmentServiceApiInstance } from 'api/managment-service';
import {
  ReactNode,
  createContext,
  useContext,
  useEffect,
  useState,
} from 'react';

export const getAuthInfo = async (): Promise<AuthenticationInfo | null> => {
  try {
    const tokenInfo = await ManagmentServiceApiInstance.getTokenInfo();
    const userInfo = await ManagmentServiceApiInstance.getUserInfo();

    return {
      expTime: tokenInfo.expTime,
      nickname: userInfo.nickname,
      email: userInfo.email,
    };
  } catch (err: unknown) {
    console.error('Failed to get user data');
    return null;
  }
};

export interface AuthenticationInfo {
  email: string;
  nickname: string;
  expTime: number;
}

export interface AuthenticationContext {
  authenticationInfo: AuthenticationInfo | null;
  setAuthenticationInfo: (value: AuthenticationInfo | null) => void;
  isTokenValid: () => Promise<boolean>;
}

export const AuthContext = createContext<AuthenticationContext>({
  authenticationInfo: null,
  setAuthenticationInfo: () => {},
  isTokenValid: () => Promise.resolve(false),
});

export const AuthProvider = (props: {
  children: ReactNode;
  authenticationInfo: AuthenticationInfo | null;
}) => {
  const [authenticationInfo, setAuthenticationInfo] =
    useState<AuthenticationInfo | null>(props.authenticationInfo);

  const isTokenValid = async (): Promise<boolean> => {
    if (authenticationInfo && authenticationInfo.expTime <= 0) {
      try {
        const authInfo = await getAuthInfo();
        setAuthenticationInfo(authInfo);
        return !!(authInfo && authInfo.expTime > 0);
      } catch {
        return false;
      }
    }

    return !!(authenticationInfo && authenticationInfo.expTime > 0);
  };

  useEffect(() => {
    setAuthenticationInfo(props.authenticationInfo);
  }, [setAuthenticationInfo]);

  return (
    <AuthContext.Provider
      value={{ authenticationInfo, setAuthenticationInfo, isTokenValid }}
    >
      {props.children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  return useContext(AuthContext);
};
