import './KindTag.scss';
import kubernetesDeployLogo from 'assets/kubernetes-deploy-icon.svg';
import kubernetesDSLogo from 'assets/kubernetes-ds-icon.svg';
import kubernetesSTSLogo from 'assets/kubernetes-sts-icon.svg';
import kubernetesLogo from 'assets/kubernetes-logo-icon.svg';

interface KindTagProps {
    name?: string;
}

const KindTag = ({ name }: KindTagProps) => {
    if (!name) {
        return <div className="kind-tag--spacer" />;
    }

    const logoMap: Record<string, string> = {
        Deployment: kubernetesDeployLogo,
        DaemonSet: kubernetesDSLogo,
        StatefulSet: kubernetesSTSLogo,
    };

    const selectedLogo = logoMap[name] || kubernetesLogo;

    return (
        <div className="kind-tag">
            <img src={selectedLogo} width="24px" height="24px" className="kind-tag__logo" />
            <div className="kind-tag__name">{name}</div>
        </div>
    );
};

export default KindTag;
