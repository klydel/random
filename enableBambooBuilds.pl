#!/usr/bin/perl -w
use WWW::Mechanize;
use HTML::TokeParser;
my $timeout = 30;
$bamuser= '';
$bampass= '';
$bamurl= 'http://server.company.com/bamboo'
sub enableAllBuilds {
    my $agent = WWW::Mechanize->new();
    $agent->get($bamurl);
    $agent->field("os_username", $bamuser);
    $agent->field("os_password", $bampass);
    $agent->click();
    $agent->get($bamurl . "/ajax/displayAllBuildSummaries.action");
    $string = $agent->content();
    while($string =~ m/resumeBuild\.action\?returnUrl(.*)buildKey(.*?)\"(.*)/g)
{
    $enableurl = $bamurl . "/build/admin/resumeBuild.action?returnUrl=%2Fbamboo%2Fstart.action&buildKey$2";
    $agent->get($enableurl);
}
}
enableAllBuilds();
