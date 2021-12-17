import json
import time
from termcolor import colored

class Issue:
    def __init__(self, issue_id, url, title, ctime):
        self.issue_id = issue_id
        self.url = url
        self.title = title
        self.ctime = ctime

    def __str__(self):
        time_format = "%d.%m.%y"
        time_str = time.strftime(time_format, self.ctime)
        return "{}: {} => {}".format(
                time_str, 
                colored(self.title, 'yellow', attrs=['bold']), 
                colored(self.url, 'blue')
                )


def make_issue(result):
    ctime_str = result["created_at"]
    time_format = "%Y-%m-%dT%H:%M:%SZ"
    ctime = time.strptime(ctime_str, time_format)
    return Issue(result["id"], result["html_url"], result["title"], ctime)


def important_issue_criteria(issue):
    return "rust-lang" not in issue.url and "deadlock" in issue.title.lower()


if __name__ == "__main__":
    filenames = ["data/deadlock4.json", "data/deadlock3.json"]
    results = []
    done_ids = {}
    for name in filenames:
        with open(name) as f:
            items = json.load(f)
            for item in items:
                item_id = item["id"]
                if not item_id in done_ids:
                    done_ids[item_id] = True
                    results.append(item)
    print("{} results found".format(colored(str(len(results)), 'red', attrs=['bold'])))
    all_issues = map(make_issue, results)
    imp_issues = list(filter(important_issue_criteria, all_issues))
    imp_issues.sort(key=lambda issue: issue.ctime, reverse=True)
    print("{} important issues found".format(colored(str(len(imp_issues)), 'red', attrs=['bold'])))
    for issue in imp_issues:
        print(str(issue))
