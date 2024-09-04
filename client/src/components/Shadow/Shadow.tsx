import "./Shadow.scss";
interface ShadowProps {
  children?: React.ReactNode;
}

const Shadow = ({ children }: ShadowProps) => {
  return <div className="shadow">{children}</div>;
};

export default Shadow;
