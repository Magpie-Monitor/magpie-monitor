import './KindTag.scss';
import kubernetesDeployLogo from 'assets/kubernetes-deploy-icon.svg';
import kubernetesDSLogo from 'assets/kubernetes-ds-icon.svg';
import kubernetesSTSLogo from 'assets/kubernetes-sts-icon.svg';
import kubernetesLogo from 'assets/kubernetes-logo-icon.svg';
import kubernetesNode from 'assets/kubernetes-node-icon.svg';
import { KUBERNETES_LINKS } from 'links/kubernetesLinks';

interface KindTagProps {
    name: string;
}

const KindTag = ({ name }: KindTagProps) => {
    const logoMap: Record<string, string> = {
        Deployment: kubernetesDeployLogo,
        DaemonSet: kubernetesDSLogo,
        StatefulSet: kubernetesSTSLogo,
        Node: kubernetesNode,
    };

    const selectedLogo = logoMap[name] || kubernetesLogo;
    const link = KUBERNETES_LINKS[name] || '#';

    return (
      <a href={link} target="_blank" rel="noopener noreferrer" className="kind-tag">
          <img src={selectedLogo} width="24px" height="24px" className="kind-tag__logo" />
          <div className="kind-tag__name">{name}</div>
      </a>
    );
};

export default KindTag;
