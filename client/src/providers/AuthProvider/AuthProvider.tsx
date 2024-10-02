import { ManagmentServiceApiInstance } from 'api/managment-service';
import { ReactNode, createContext, useContext, useEffect, useState } from 'react';

export const getAuthInfo = async (): Promise<AuthenticationInfo> => {
  const tokenInfo = await ManagmentServiceApiInstance.getTokenInfo();
  const userInfo = await ManagmentServiceApiInstance.getUserInfo();

  return {
    expTime: tokenInfo.expTime,
    nickname: userInfo.nickname,
    email: userInfo.email,
  };
};

export interface AuthenticationInfo {
  email: string;
  nickname: string;
  expTime: number;
}

export interface AuthenticationContext {
  authenticationInfo: AuthenticationInfo;
  setAuthenticationInfo: (value: AuthenticationInfo) => void;
  isTokenValid: () => Promise<boolean>;
}

export const AuthContext = createContext<AuthenticationContext>({
  authenticationInfo: { email: '', expTime: 0, nickname: '' },
  setAuthenticationInfo: () => {},
  isTokenValid: () => Promise.resolve(false),
});

export const AuthProvider = (props: {
  children: ReactNode;
  authenticationInfo: AuthenticationInfo;
}) => {
  const [authenticationInfo, setAuthenticationInfo] = useState<AuthenticationInfo>(
    props.authenticationInfo,
  );

  const isTokenValid = async (): Promise<boolean> => {
    if (authenticationInfo.email && authenticationInfo.expTime <= 0) {
      try {
        const authInfo = await getAuthInfo();
        setAuthenticationInfo(authInfo);
        return !!(authInfo.expTime && authInfo.expTime > 0);
      } catch {
        return false;
      }
    }

    return !!(authenticationInfo.expTime && authenticationInfo.expTime > 0);
  };

  useEffect(() => {
    setAuthenticationInfo(props.authenticationInfo);
  }, [props, setAuthenticationInfo]);

  return (
    <AuthContext.Provider value={{ authenticationInfo, setAuthenticationInfo, isTokenValid }}>
      {props.children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  return useContext(AuthContext);
};
