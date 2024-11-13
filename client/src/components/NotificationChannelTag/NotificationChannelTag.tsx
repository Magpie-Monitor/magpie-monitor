import BrandTag from 'components/BrandTag/BrandTag.tsx';

interface NotificationChannelTagProps {
  name: string;
  logoImgSrc: string;
}

const NotificationChannelTag = ({
  name,
  logoImgSrc,
}: NotificationChannelTagProps) => {
  return (
    <BrandTag
      name={name}
      logo={<img src={logoImgSrc} width="20px" height="20px" />}
    />
  );
};

export default NotificationChannelTag;
