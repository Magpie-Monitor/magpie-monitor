import { useNavigate, useOutlet } from 'react-router-dom';
import { useAuth } from './AuthProvider';
import { useEffect } from 'react';
import Navbar from 'components/Navbar/Navbar';
import './ProtectedLayout.scss';

export const ProtectedLayout = () => {
  const { isTokenValid } = useAuth();
  const outlet = useOutlet();
  const navigate = useNavigate();

  useEffect(() => {
    isTokenValid().then((isValid) => {
      if (!isValid) {
        navigate('/login');
      }
    });
  }, [isTokenValid, navigate]);

  return (
    <div className="protected-layout">
      <Navbar />
      {outlet}
    </div>
  );
};
