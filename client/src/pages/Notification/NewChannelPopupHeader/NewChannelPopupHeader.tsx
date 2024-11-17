import './NewChannelPopupHeader.scss';

interface NewChannelPopupHeaderProps {
  icon?: React.ReactNode;
  title: React.ReactNode;
}

const NewChannelPopupHeader = ({ icon, title }: NewChannelPopupHeaderProps) => {
  return (
    <div className="new-channel-popup-header ">
      <div className="new-channel-popup-header__icon">{icon}</div>
      <div className="new-channel-popup-header__title">{title}</div>
    </div>
  );
};

export default NewChannelPopupHeader;
