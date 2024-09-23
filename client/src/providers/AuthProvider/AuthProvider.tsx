import {
  ReactNode,
  createContext,
  useContext,
  useEffect,
  useState,
} from "react";

export interface AuthenticationInfo {
  email: string;
  expires: number;
}

export interface AuthenticationContext {
  authenticationInfo: AuthenticationInfo;
  setAuthenticationInfo: (value: AuthenticationInfo) => void;
  isTokenValid: () => boolean;
}

export const AuthContext = createContext<AuthenticationContext>({
  authenticationInfo: { email: "", expires: 0 },
  setAuthenticationInfo: () => { },
  isTokenValid: () => false,
});

export const AuthProvider = (props: {
  children: ReactNode;
  authenticationInfo: AuthenticationInfo;
}) => {
  const [authenticationInfo, setAuthenticationInfo] =
    useState<AuthenticationInfo>(props.authenticationInfo);

  const isTokenValid = () => {
    return !!(
      authenticationInfo.expires &&
      authenticationInfo.expires > Date.now() / 1000
    );
  };

  useEffect(() => {
    setAuthenticationInfo(props.authenticationInfo);
  }, [props, setAuthenticationInfo]);

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
