import './HeaderWithIcon.scss';

interface HeaderWithIconProps {
  icon?: React.ReactNode;
  title: React.ReactNode;
}

const HeaderWithIcon = ({ icon, title }: HeaderWithIconProps) => {
  return (
    <div className="header-with-icon">
      <div className="header-with-icon__icon">{icon}</div>
      <div className="header-with-icon__title">{title}</div>
    </div>
  );
};

export default HeaderWithIcon;
