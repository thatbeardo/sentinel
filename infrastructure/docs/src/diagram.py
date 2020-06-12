from diagrams import Cluster, Diagram, Edge
from diagrams.custom import Custom
from diagrams.aws.compute import EC2, ECS
from diagrams.aws.network import CF, ELB, Route53
from diagrams.aws.storage import S3

from urllib.request import urlretrieve

with Diagram(
    "GuardMy.App AWS Infrastructure", 
    direction="TB",
    filename="guard_my_app_infrastructure",
    outformat="png",
):
    dns = Route53("*.guardmy.app")

    load_balancer = ELB("api.guardmy.app")
    
    with Cluster("Sentinel API"):
        with Cluster("Fargate"):
            sentinel_api = [
                ECS("task 1"),
                ECS("task 2"),
                ECS("task 3")
            ]
        
        neo4j = EC2("Neo4j")

    dns >> load_balancer >> sentinel_api >> neo4j


    cloudfront = CF('cdn')
    ui = S3("www.guardmy.app")
    auth0 = Custom("Auth0", './assets/auth0.png')
    dns >> cloudfront >> ui >> auth0 >> ui >> load_balancer

