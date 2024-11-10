import CustomTag from 'components/BrandTag/CustomTag.tsx';

interface NotificationChannelTagProps {
  name: string;
  logoImgSrc: string;
}

const NotificationChannelTag = ({
  name,
  logoImgSrc,
}: NotificationChannelTagProps) => {
  return (
    <CustomTag
      name={name}
      logo={<img src={logoImgSrc} width="20px" height="20px" />}
    />
  );
};

export default NotificationChannelTag;
