import './ClusterLink.scss';

interface ClusterLinkParams {
  name: string;
}

const ClusterLink = ({ name }: ClusterLinkParams) => {
  return <div className={'cluster-link'}>{name}</div>;
};

export default ClusterLink;
